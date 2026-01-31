package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/internal/db"
	"server/internal/models"
	"server/internal/response"
	"server/pkg/config"

	"github.com/gin-gonic/gin"
)

func TestGetDashboardStats_FromOrders(t *testing.T) {
	testDB := setupOrderTestDB(t)
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	originalLoc := config.GlobalBusinessHours.TimeLocation
	config.GlobalBusinessHours.TimeLocation = time.UTC
	defer func() { config.GlobalBusinessHours.TimeLocation = originalLoc }()

	now := time.Now().In(time.UTC)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	memberToday := models.Member{Name: "A", Phone: "10000000041", InvitationCode: "code-10000000041"}
	memberMonth := models.Member{Name: "B", Phone: "10000000042", InvitationCode: "code-10000000042"}
	memberLastMonth := models.Member{Name: "C", Phone: "10000000043", InvitationCode: "code-10000000043"}
	testDB.Create(&memberToday)
	testDB.Create(&memberMonth)
	testDB.Create(&memberLastMonth)

	testDB.Model(&models.Member{}).Where("id = ?", memberToday.ID).Updates(map[string]interface{}{
		"created_at": today.Add(2 * time.Hour),
		"updated_at": today.Add(2 * time.Hour),
	})
	testDB.Model(&models.Member{}).Where("id = ?", memberMonth.ID).Updates(map[string]interface{}{
		"created_at": monthStart.Add(24 * time.Hour),
		"updated_at": monthStart.Add(24 * time.Hour),
	})
	testDB.Model(&models.Member{}).Where("id = ?", memberLastMonth.ID).Updates(map[string]interface{}{
		"created_at": monthStart.AddDate(0, -1, 0).Add(24 * time.Hour),
		"updated_at": monthStart.AddDate(0, -1, 0).Add(24 * time.Hour),
	})

	tech := models.Technician{Name: "T", Status: 0}
	service := models.ServiceProduct{Name: "S", Duration: 60, Price: 100}
	product := models.PhysicalProduct{Name: "P", Stock: 100, RetailPrice: 20, CostPrice: 10, IsActive: true}
	operator := models.User{Username: "op-stats", PasswordHash: "x", Role: "operator", IsActive: true}
	testDB.Create(&tech)
	testDB.Create(&service)
	testDB.Create(&product)
	testDB.Create(&operator)

	appt := models.Appointment{
		MemberID:    memberToday.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   today.Add(9 * time.Hour),
		EndTime:     today.Add(10 * time.Hour),
		Status:      "completed",
		OriginPrice: 100,
		ActualPrice: 80,
	}
	testDB.Create(&appt)

	apptYesterday := models.Appointment{
		MemberID:    memberToday.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   yesterday.Add(9 * time.Hour),
		EndTime:     yesterday.Add(10 * time.Hour),
		Status:      "completed",
		OriginPrice: 100,
		ActualPrice: 80,
	}
	testDB.Create(&apptYesterday)

	saleAmount := 40.0
	invLog := models.InventoryLog{
		ProductID:    product.ID,
		OperatorID:   operator.ID,
		MemberID:     &memberToday.ID,
		ChangeAmount: -2,
		ActionType:   "sale",
		BeforeStock:  100,
		AfterStock:   98,
		SaleAmount:   &saleAmount,
	}
	testDB.Create(&invLog)

	orderServiceToday := models.Order{MemberID: memberToday.ID, PaidAmount: 30, CommissionAmount: 0, OrderType: "service", AppointmentID: &appt.ID}
	orderPhysicalToday := models.Order{MemberID: memberToday.ID, PaidAmount: 70, CommissionAmount: 0, OrderType: "physical", InventoryLogID: &invLog.ID}
	orderYesterday := models.Order{MemberID: memberToday.ID, PaidAmount: 50, CommissionAmount: 0, OrderType: "service", AppointmentID: &apptYesterday.ID}
	testDB.Create(&orderServiceToday)
	testDB.Create(&orderPhysicalToday)
	testDB.Create(&orderYesterday)

	testDB.Model(&models.Order{}).Where("id = ?", orderServiceToday.ID).Updates(map[string]interface{}{
		"created_at": today.Add(1 * time.Hour),
		"updated_at": today.Add(1 * time.Hour),
	})
	testDB.Model(&models.Order{}).Where("id = ?", orderPhysicalToday.ID).Updates(map[string]interface{}{
		"created_at": today.Add(3 * time.Hour),
		"updated_at": today.Add(3 * time.Hour),
	})
	testDB.Model(&models.Order{}).Where("id = ?", orderYesterday.ID).Updates(map[string]interface{}{
		"created_at": yesterday.Add(2 * time.Hour),
		"updated_at": yesterday.Add(2 * time.Hour),
	})

	pendingAppt := models.Appointment{
		MemberID:    memberToday.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   today.Add(11 * time.Hour),
		EndTime:     today.Add(12 * time.Hour),
		Status:      "pending",
		OriginPrice: 100,
		ActualPrice: 80,
	}
	waitingAppt := models.Appointment{
		MemberID:    memberToday.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   today.Add(11 * time.Hour),
		EndTime:     today.Add(12 * time.Hour),
		Status:      "waiting",
		OriginPrice: 100,
		ActualPrice: 80,
	}
	testDB.Create(&pendingAppt)
	testDB.Create(&waitingAppt)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/dashboard/stats", GetDashboardStats)

	req, _ := http.NewRequest("GET", "/api/dashboard/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var env response.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	dataBytes, _ := json.Marshal(env.Data)
	var data struct {
		DailyRevenue         float64 `json:"dailyRevenue"`
		RevenueGrowth        float64 `json:"revenueGrowth"`
		NewMembers           int64   `json:"newMembers"`
		MonthlyNewMembers    int64   `json:"monthlyNewMembers"`
		PendingAppointments  int64   `json:"pendingAppointments"`
	}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}

	if data.DailyRevenue != 100 {
		t.Fatalf("expected dailyRevenue=100, got %v", data.DailyRevenue)
	}
	if data.RevenueGrowth < 99.9 || data.RevenueGrowth > 100.1 {
		t.Fatalf("expected revenueGrowth≈100, got %v", data.RevenueGrowth)
	}
	if data.NewMembers != 1 {
		t.Fatalf("expected newMembers=1, got %d", data.NewMembers)
	}
	if data.MonthlyNewMembers != 2 {
		t.Fatalf("expected monthlyNewMembers=2, got %d", data.MonthlyNewMembers)
	}
	if data.PendingAppointments != 2 {
		t.Fatalf("expected pendingAppointments=2, got %d", data.PendingAppointments)
	}
}
