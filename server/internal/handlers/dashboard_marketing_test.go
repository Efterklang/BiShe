package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/internal/models"
	"server/internal/response"

	"github.com/gin-gonic/gin"
)

func TestGetMarketingMetrics_SummaryAndFilters(t *testing.T) {
	testDB := setupOrderTestDB(t)

	m1 := models.Member{Name: "M1", Phone: "10000000011", Level: "basic", InvitationCode: "code-10000000011"}
	m2 := models.Member{Name: "M2", Phone: "10000000012", Level: "gold", InvitationCode: "code-10000000012"}
	if err := testDB.Create(&m1).Error; err != nil {
		t.Fatalf("create member1: %v", err)
	}
	if err := testDB.Create(&m2).Error; err != nil {
		t.Fatalf("create member2: %v", err)
	}

	operator := models.User{Username: "op2", PasswordHash: "x", Role: "operator", IsActive: true}
	product := models.PhysicalProduct{Name: "Oil", Stock: 10, RetailPrice: 20, CostPrice: 10, IsActive: true}
	testDB.Create(&operator)
	testDB.Create(&product)

	saleAmount1 := 40.0
	saleAmount2 := 20.0

	log1 := models.InventoryLog{ProductID: product.ID, OperatorID: operator.ID, MemberID: &m1.ID, ChangeAmount: -2, ActionType: "sale", BeforeStock: 10, AfterStock: 8, SaleAmount: &saleAmount1}
	log2 := models.InventoryLog{ProductID: product.ID, OperatorID: operator.ID, MemberID: &m1.ID, ChangeAmount: -1, ActionType: "sale", BeforeStock: 8, AfterStock: 7, SaleAmount: &saleAmount2}
	testDB.Create(&log1)
	testDB.Create(&log2)

	tech := models.Technician{Name: "Bob", Status: 0}
	service := models.ServiceProduct{Name: "Massage", Duration: 60, Price: 100}
	testDB.Create(&tech)
	testDB.Create(&service)

	t1 := time.Date(2026, 1, 1, 10, 0, 0, 0, time.Local)
	t2 := time.Date(2026, 1, 2, 10, 0, 0, 0, time.Local)

	appt := models.Appointment{
		MemberID:    m2.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   time.Now().Add(-2 * time.Hour),
		EndTime:     time.Now().Add(-1 * time.Hour),
		Status:      "completed",
		OriginPrice: 100,
		ActualPrice: 50,
	}
	appt.CreatedAt = t1
	appt.UpdatedAt = t1
	testDB.Create(&appt)
	testDB.Model(&models.Appointment{}).Where("id = ?", appt.ID).Update("created_at", t1)
	testDB.Model(&models.Appointment{}).Where("id = ?", appt.ID).Update("updated_at", t1)

	o1 := models.Order{MemberID: m1.ID, PaidAmount: 40, CommissionAmount: 4, OrderType: "physical", InventoryLogID: &log1.ID}
	o1.CreatedAt = t1
	o1.UpdatedAt = t1
	o2 := models.Order{MemberID: m1.ID, PaidAmount: 20, CommissionAmount: 2, OrderType: "physical", InventoryLogID: &log2.ID}
	o2.CreatedAt = t2
	o2.UpdatedAt = t2
	o3 := models.Order{MemberID: m2.ID, PaidAmount: 50, CommissionAmount: 5, OrderType: "service", AppointmentID: &appt.ID}
	o3.CreatedAt = t1
	o3.UpdatedAt = t1
	testDB.Create(&o1)
	testDB.Create(&o2)
	testDB.Create(&o3)

	var apptCount int64
	if err := testDB.Table("appointments").Where("substr(created_at, 1, 10) >= ? AND substr(created_at, 1, 10) <= ?", "2026-01-01", "2026-01-02").Count(&apptCount).Error; err != nil {
		t.Fatalf("count appointments: %v", err)
	}
	if apptCount != 1 {
		t.Fatalf("expected appointment count=1 in range, got %d", apptCount)
	}

	handler := NewDashboardHandler(testDB)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/dashboard/marketing", handler.GetMarketingMetrics)

	req, _ := http.NewRequest("GET", "/api/dashboard/marketing?start=2026-01-01&end=2026-01-02&granularity=day", nil)
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
		Summary struct {
			TotalSales     float64 `json:"total_sales"`
			BuyerCount     int64   `json:"buyer_count"`
			RepurchaseRate float64 `json:"repurchase_rate"`
			ConversionRate float64 `json:"conversion_rate"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}

	if data.Summary.TotalSales != 110 {
		t.Fatalf("expected total_sales=110, got %v", data.Summary.TotalSales)
	}
	if data.Summary.BuyerCount != 2 {
		t.Fatalf("expected buyer_count=2, got %d", data.Summary.BuyerCount)
	}
	if data.Summary.RepurchaseRate < 49.9 || data.Summary.RepurchaseRate > 50.1 {
		t.Fatalf("expected repurchase_rate≈50, got %v", data.Summary.RepurchaseRate)
	}
	if data.Summary.ConversionRate < 99.9 || data.Summary.ConversionRate > 100.1 {
		t.Fatalf("expected conversion_rate≈100, got %v", data.Summary.ConversionRate)
	}

	req2, _ := http.NewRequest("GET", "/api/dashboard/marketing?start=2026-01-01&end=2026-01-02&member_level=basic", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w2.Code, w2.Body.String())
	}

	var env2 response.Envelope
	json.Unmarshal(w2.Body.Bytes(), &env2)
	dataBytes2, _ := json.Marshal(env2.Data)
	var data2 struct {
		Summary struct {
			BuyerCount     int64   `json:"buyer_count"`
			RepurchaseRate float64 `json:"repurchase_rate"`
		} `json:"summary"`
	}
	json.Unmarshal(dataBytes2, &data2)
	if data2.Summary.BuyerCount != 1 {
		t.Fatalf("expected buyer_count=1, got %d", data2.Summary.BuyerCount)
	}
	if data2.Summary.RepurchaseRate < 99.9 || data2.Summary.RepurchaseRate > 100.1 {
		t.Fatalf("expected repurchase_rate≈100, got %v", data2.Summary.RepurchaseRate)
	}
}
