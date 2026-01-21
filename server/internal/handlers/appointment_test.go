package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestListAppointments_PreloadFix(t *testing.T) {
	// 1. Setup in-memory DB
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}
	// Migrate all related models
	err = testDB.AutoMigrate(
		&models.Member{},
		&models.Technician{},
		&models.ServiceProduct{},
		&models.Appointment{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Swap global DB
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// 2. Insert dummy data to ensure Preload has something to load
	member := models.Member{Name: "Alice", Phone: "1234567890"}
	testDB.Create(&member)

	tech := models.Technician{Name: "Bob", Status: 0}
	testDB.Create(&tech)

	service := models.ServiceProduct{Name: "Massage", Duration: 60, Price: 100}
	testDB.Create(&service)

	appt := models.Appointment{
		MemberID:    member.ID,
		TechID:      tech.ID,
		ServiceID:   service.ID,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Status:      "pending",
		OriginPrice: 100,
		ActualPrice: 100,
	}
	if err := testDB.Create(&appt).Error; err != nil {
		t.Fatalf("Failed to create appointment: %v", err)
	}

	// 3. Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/appointments", ListAppointments)

	// 4. Make Request
	req, _ := http.NewRequest("GET", "/appointments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 5. Assertions
	// Before fix, this should fail with 400 and error message about unsupported relations
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify that the response contains the preloaded data
	body := w.Body.String()
	if !strings.Contains(body, "Massage") {
		t.Errorf("Response body should contain service name 'Massage', got: %s", body)
	}
}
