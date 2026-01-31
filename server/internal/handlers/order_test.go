package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"server/internal/db"
	"server/internal/models"
	"server/pkg/config"
	"server/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var orderTestDBSeq uint64

func setupOrderTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	seq := atomic.AddUint64(&orderTestDBSeq, 1)
	dsn := fmt.Sprintf("file:orders_test_%d?mode=memory&cache=shared", seq)
	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}
	sqlDB, err := database.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
	}
	if err := database.AutoMigrate(
		&models.User{},
		&models.Member{},
		&models.Technician{},
		&models.ServiceProduct{},
		&models.Appointment{},
		&models.PhysicalProduct{},
		&models.InventoryLog{},
		&models.Order{},
		&models.FissionLog{},
	); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	return database
}

func TestCreateOrder_ServiceFromCompletedAppointment(t *testing.T) {
	testDB := setupOrderTestDB(t)
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	member := models.Member{Name: "Alice", Phone: "10000000001", InvitationCode: "code-10000000001"}
	if err := testDB.Create(&member).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}

	tech := models.Technician{Name: "Bob", Status: 0}
	service := models.ServiceProduct{Name: "Massage", Duration: 60, Price: 100}
	testDB.Create(&tech)
	testDB.Create(&service)

	appt := models.Appointment{
		MemberID:    member.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   time.Now().Add(-2 * time.Hour),
		EndTime:     time.Now().Add(-1 * time.Hour),
		Status:      "completed",
		OriginPrice: 100,
		ActualPrice: 80.5,
	}
	if err := testDB.Create(&appt).Error; err != nil {
		t.Fatalf("create appointment: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/orders", CreateOrder)

	body, _ := json.Marshal(gin.H{
		"order_type":     "service",
		"appointment_id": appt.ID,
	})
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var count int64
	if err := testDB.Model(&models.Order{}).Count(&count).Error; err != nil {
		t.Fatalf("count orders: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 order, got %d", count)
	}
}

func TestCreateOrder_ServiceRejectPendingAppointment(t *testing.T) {
	testDB := setupOrderTestDB(t)
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	member := models.Member{Name: "Alice", Phone: "10000000002", InvitationCode: "code-10000000002"}
	testDB.Create(&member)

	tech := models.Technician{Name: "Bob", Status: 0}
	service := models.ServiceProduct{Name: "Massage", Duration: 60, Price: 100}
	testDB.Create(&tech)
	testDB.Create(&service)

	appt := models.Appointment{
		MemberID:    member.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   time.Now().Add(1 * time.Hour),
		EndTime:     time.Now().Add(2 * time.Hour),
		Status:      "pending",
		OriginPrice: 100,
		ActualPrice: 80,
	}
	testDB.Create(&appt)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/orders", CreateOrder)

	body, _ := json.Marshal(gin.H{
		"order_type":     "service",
		"appointment_id": appt.ID,
	})
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCreateOrder_PhysicalFromInventoryLog(t *testing.T) {
	testDB := setupOrderTestDB(t)
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	operator := models.User{Username: "op", PasswordHash: "x", Role: "operator", IsActive: true}
	member := models.Member{Name: "Alice", Phone: "10000000003", InvitationCode: "code-10000000003"}
	product := models.PhysicalProduct{Name: "Oil", Stock: 10, RetailPrice: 20, CostPrice: 10, IsActive: true}
	testDB.Create(&operator)
	testDB.Create(&member)
	testDB.Create(&product)

	saleAmount := 40.0
	invLog := models.InventoryLog{
		ProductID:    product.ID,
		OperatorID:   operator.ID,
		MemberID:     &member.ID,
		ChangeAmount: -2,
		ActionType:   "sale",
		BeforeStock:  10,
		AfterStock:   8,
		SaleAmount:   &saleAmount,
		Remark:       "test",
	}
	if err := testDB.Create(&invLog).Error; err != nil {
		t.Fatalf("create inventory log: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/orders", CreateOrder)

	body, _ := json.Marshal(gin.H{
		"order_type":       "physical",
		"inventory_log_id": invLog.ID,
	})
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var count int64
	testDB.Model(&models.Order{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 order, got %d", count)
	}
}

func TestOrderIntegrityConstraints_DBLevel(t *testing.T) {
	testDB := setupOrderTestDB(t)

	member := models.Member{Name: "Alice", Phone: "10000000004", InvitationCode: "code-10000000004"}
	testDB.Create(&member)

	invalid := models.Order{
		MemberID:         member.ID,
		PaidAmount:       10,
		CommissionAmount: 0,
		OrderType:        "service",
	}
	if err := testDB.Create(&invalid).Error; err == nil {
		t.Fatalf("expected constraint error, got nil")
	}
}

func TestCreateOrder_ConcurrentIdempotency(t *testing.T) {
	testDB := setupOrderTestDB(t)
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	member := models.Member{Name: "Alice", Phone: "10000000005", InvitationCode: "code-10000000005"}
	testDB.Create(&member)
	tech := models.Technician{Name: "Bob", Status: 0}
	service := models.ServiceProduct{Name: "Massage", Duration: 60, Price: 100}
	testDB.Create(&tech)
	testDB.Create(&service)

	appt := models.Appointment{
		MemberID:    member.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   time.Now().Add(-2 * time.Hour),
		EndTime:     time.Now().Add(-1 * time.Hour),
		Status:      "completed",
		OriginPrice: 100,
		ActualPrice: 99,
	}
	testDB.Create(&appt)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/orders", CreateOrder)

	body, _ := json.Marshal(gin.H{
		"order_type":     "service",
		"appointment_id": appt.ID,
	})

	var wg sync.WaitGroup
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("POST", "/api/orders", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}()
	}
	wg.Wait()

	var count int64
	testDB.Model(&models.Order{}).Where("appointment_id = ?", appt.ID).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 order for appointment, got %d", count)
	}
}

func TestCompleteAppointment_CreatesServiceOrder(t *testing.T) {
	testDB := setupOrderTestDB(t)
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	referrer := models.Member{Name: "Ref", Phone: "10000000006", InvitationCode: "code-10000000006", Balance: 0}
	invitee := models.Member{Name: "Inv", Phone: "10000000007", InvitationCode: "code-10000000007"}
	testDB.Create(&referrer)
	invitee.ReferrerID = &referrer.ID
	testDB.Create(&invitee)

	tech := models.Technician{Name: "Bob", Status: 0}
	service := models.ServiceProduct{Name: "Massage", Duration: 60, Price: 100}
	testDB.Create(&tech)
	testDB.Create(&service)

	appt := models.Appointment{
		MemberID:    invitee.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   time.Now().Add(-2 * time.Hour),
		EndTime:     time.Now().Add(-1 * time.Hour),
		Status:      "pending",
		OriginPrice: 100,
		ActualPrice: 50,
	}
	testDB.Create(&appt)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/api/appointments/:id/complete", CompleteAppointment)

	pay := map[string]interface{}{
		"payment_method": "cash",
		"balance_amount": 0,
		"cash_amount":    50,
	}
	payBody, _ := json.Marshal(pay)
	req, _ := http.NewRequest("PUT", "/api/appointments/"+strconvUint(appt.ID)+"/complete", bytes.NewReader(payBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var order models.Order
	if err := testDB.Where("appointment_id = ?", appt.ID).First(&order).Error; err != nil {
		t.Fatalf("expected service order, err=%v", err)
	}
	if order.OrderType != "service" {
		t.Fatalf("expected order_type service, got %s", order.OrderType)
	}
	if order.InviterID == nil || *order.InviterID != referrer.ID {
		t.Fatalf("expected inviter_id %d, got %v", referrer.ID, order.InviterID)
	}
	expectedCommission := util.CalculateRate(appt.ActualPrice, config.GlobalCommission.ReferralRate)
	if util.ToCents(order.CommissionAmount) != util.ToCents(expectedCommission) {
		t.Fatalf("expected commission %.2f, got %.2f", expectedCommission, order.CommissionAmount)
	}
}

func strconvUint(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
