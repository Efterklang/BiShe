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

func TestDashboardRevenueTrend_FromOrders(t *testing.T) {
	testDB := setupOrderTestDB(t)

	m := models.Member{Name: "M", Phone: "10000000031", InvitationCode: "code-10000000031"}
	if err := testDB.Create(&m).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}

	now := time.Now()
	d1 := now.AddDate(0, 0, -1)

	serviceOrder := models.Order{MemberID: m.ID, PaidAmount: 10, CommissionAmount: 0, OrderType: "service", AppointmentID: ptrUint(1)}
	serviceOrder.CreatedAt = d1
	serviceOrder.UpdatedAt = d1
	physicalOrder := models.Order{MemberID: m.ID, PaidAmount: 7, CommissionAmount: 0, OrderType: "physical", InventoryLogID: ptrUint(1)}
	physicalOrder.CreatedAt = d1
	physicalOrder.UpdatedAt = d1
	if err := testDB.Create(&serviceOrder).Error; err != nil {
		t.Fatalf("create service order: %v", err)
	}
	if err := testDB.Create(&physicalOrder).Error; err != nil {
		t.Fatalf("create physical order: %v", err)
	}

	handler := NewDashboardHandler(testDB)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/dashboard/revenue-trend", handler.GetRevenueTrend)

	req, _ := http.NewRequest("GET", "/api/dashboard/revenue-trend?days=3", nil)
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
	var rows []struct {
		Date           string  `json:"date"`
		ServiceRevenue float64 `json:"service_revenue"`
		ProductRevenue float64 `json:"product_revenue"`
	}
	if err := json.Unmarshal(dataBytes, &rows); err != nil {
		t.Fatalf("unmarshal rows: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 days, got %d", len(rows))
	}

	target := d1.Format("2006-01-02")
	found := false
	for _, r := range rows {
		if r.Date == target {
			found = true
			if r.ServiceRevenue != 10 {
				t.Fatalf("expected service_revenue=10, got %v", r.ServiceRevenue)
			}
			if r.ProductRevenue != 7 {
				t.Fatalf("expected product_revenue=7, got %v", r.ProductRevenue)
			}
		}
	}
	if !found {
		t.Fatalf("expected to find date %s", target)
	}
}

func TestDashboardServiceRanking_FromOrders(t *testing.T) {
	testDB := setupOrderTestDB(t)

	member := models.Member{Name: "M", Phone: "10000000032", InvitationCode: "code-10000000032"}
	tech := models.Technician{Name: "T", Status: 0}
	service := models.ServiceProduct{Name: "S", Duration: 60, Price: 100}
	if err := testDB.Create(&member).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}
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
		ActualPrice: 80,
	}
	if err := testDB.Create(&appt).Error; err != nil {
		t.Fatalf("create appointment: %v", err)
	}

	order := models.Order{MemberID: member.ID, PaidAmount: 80, CommissionAmount: 0, OrderType: "service", AppointmentID: &appt.ID}
	if err := testDB.Create(&order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	handler := NewDashboardHandler(testDB)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/dashboard/service-ranking", handler.GetServiceRanking)

	req, _ := http.NewRequest("GET", "/api/dashboard/service-ranking", nil)
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
	var ranks []struct {
		ServiceID    uint    `json:"service_id"`
		ServiceName  string  `json:"service_name"`
		OrderCount   int64   `json:"order_count"`
		TotalRevenue float64 `json:"total_revenue"`
	}
	if err := json.Unmarshal(dataBytes, &ranks); err != nil {
		t.Fatalf("unmarshal ranks: %v", err)
	}
	if len(ranks) == 0 {
		t.Fatalf("expected non-empty ranks")
	}
	if ranks[0].ServiceID != service.ID || ranks[0].OrderCount != 1 || ranks[0].TotalRevenue != 80 {
		t.Fatalf("unexpected rank: %+v", ranks[0])
	}
}

func TestDashboardProductSalesOverview_FromOrders(t *testing.T) {
	testDB := setupOrderTestDB(t)

	operator := models.User{Username: "opx", PasswordHash: "x", Role: "operator", IsActive: true}
	member := models.Member{Name: "M", Phone: "10000000033", InvitationCode: "code-10000000033"}
	product := models.PhysicalProduct{Name: "P", Stock: 100, RetailPrice: 20, CostPrice: 10, IsActive: true}
	testDB.Create(&operator)
	testDB.Create(&member)
	testDB.Create(&product)

	log := models.InventoryLog{
		ProductID:    product.ID,
		OperatorID:   operator.ID,
		MemberID:     &member.ID,
		ChangeAmount: -2,
		ActionType:   "sale",
		BeforeStock:  100,
		AfterStock:   98,
		Remark:       "x",
	}
	if err := testDB.Create(&log).Error; err != nil {
		t.Fatalf("create inventory log: %v", err)
	}

	order := models.Order{MemberID: member.ID, PaidAmount: 40, CommissionAmount: 0, OrderType: "physical", InventoryLogID: &log.ID}
	if err := testDB.Create(&order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	handler := NewDashboardHandler(testDB)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/dashboard/product-sales", handler.GetProductSalesOverview)

	req, _ := http.NewRequest("GET", "/api/dashboard/product-sales?days=30", nil)
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
		TopProducts []struct {
			ProductID    uint    `json:"product_id"`
			SalesCount   int64   `json:"sales_count"`
			TotalRevenue float64 `json:"total_revenue"`
		} `json:"topProducts"`
		TotalRevenue float64 `json:"totalRevenue"`
		TotalSales   int64   `json:"totalSales"`
	}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}

	if data.TotalRevenue != 40 {
		t.Fatalf("expected totalRevenue=40, got %v", data.TotalRevenue)
	}
	if data.TotalSales != 1 {
		t.Fatalf("expected totalSales=1, got %d", data.TotalSales)
	}
	if len(data.TopProducts) == 0 || data.TopProducts[0].ProductID != product.ID {
		t.Fatalf("expected top product %d, got %+v", product.ID, data.TopProducts)
	}
}

func ptrUint(v uint) *uint {
	return &v
}

