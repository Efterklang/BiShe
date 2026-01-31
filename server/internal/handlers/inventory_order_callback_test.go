package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/db"
	"server/internal/models"

	"github.com/gin-gonic/gin"
)

func TestCreateInventoryChange_SaleCreatesPhysicalOrder(t *testing.T) {
	testDB := setupOrderTestDB(t)
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	operator := models.User{Username: "op3", PasswordHash: "x", Role: "operator", IsActive: true}
	member := models.Member{Name: "Alice", Phone: "10000000021", InvitationCode: "code-10000000021"}
	product := models.PhysicalProduct{Name: "Oil", Stock: 10, RetailPrice: 20, CostPrice: 10, IsActive: true}
	testDB.Create(&operator)
	testDB.Create(&member)
	testDB.Create(&product)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/inventory/change", func(c *gin.Context) {
		c.Set("user_id", operator.ID)
		CreateInventoryChange(c)
	})

	reqBody, _ := json.Marshal(gin.H{
		"product_id":     product.ID,
		"change_amount":  -2,
		"action_type":    "sale",
		"member_id":      member.ID,
		"remark":         "test",
	})
	req, _ := http.NewRequest("POST", "/api/inventory/change", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var orderCount int64
	testDB.Model(&models.Order{}).Count(&orderCount)
	if orderCount != 1 {
		t.Fatalf("expected 1 order, got %d", orderCount)
	}

	var order models.Order
	if err := testDB.First(&order).Error; err != nil {
		t.Fatalf("read order: %v", err)
	}
	if order.OrderType != "physical" {
		t.Fatalf("expected order_type physical, got %s", order.OrderType)
	}
	if order.MemberID != member.ID {
		t.Fatalf("expected member_id %d, got %d", member.ID, order.MemberID)
	}
}
