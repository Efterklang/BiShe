package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"

	// "smartspa-admin/internal/repo" // Implicitly used by handlers

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFilterTestDB() *gorm.DB {
	d, err := gorm.Open(sqlite.Open("file:filtertest?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	d.AutoMigrate(&models.Technician{}, &models.ServiceProduct{}, &models.Schedule{}, &models.Appointment{})
	return d
}

func TestGetAvailableTechnicians_SkillFilter(t *testing.T) {
	// Setup DB
	testDB := setupFilterTestDB()
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/api/schedules/available-technicians", GetAvailableTechnicians)

	// Create Services
	massage := models.ServiceProduct{BaseModel: models.BaseModel{ID: 1}, Name: "Massage", Duration: 60, Price: 100}
	facial := models.ServiceProduct{BaseModel: models.BaseModel{ID: 2}, Name: "Facial", Duration: 60, Price: 100}
	testDB.Create(&massage)
	testDB.Create(&facial)

	// Create Technicians with Skills
	// Tech 1: Has Massage (New Format: [1])
	skills1, _ := json.Marshal([]uint{1})
	testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: 1}, Name: "Tech 1", Skills: skills1})

	// Tech 2: Has Facial (New Format: [2])
	skills2, _ := json.Marshal([]uint{2})
	testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: 2}, Name: "Tech 2", Skills: skills2})

	// Tech 3: Has Massage (Old Format: ["Massage"])
	skills3, _ := json.Marshal([]string{"Massage"})
	testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: 3}, Name: "Tech 3", Skills: skills3})

	// Tech 4: No Skills
	testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: 4}, Name: "Tech 4"})

	// Test Case: Filter for Massage (Service ID 1)
	t.Run("Filter by Skill Massage", func(t *testing.T) {
		w := httptest.NewRecorder()
		startTime := time.Now().Add(time.Hour).Format(time.RFC3339)

		params := url.Values{}
		params.Add("service_id", "1")
		params.Add("start_time", startTime)

		req, _ := http.NewRequest("GET", "/api/schedules/available-technicians?"+params.Encode(), nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Data struct {
				Available   []models.Technician `json:"available"`
				Unavailable []models.Technician `json:"unavailable"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Expected Available: Tech 1, Tech 3
		availableIDs := make(map[uint]bool)
		for _, t := range response.Data.Available {
			availableIDs[t.ID] = true
		}

		if !availableIDs[1] {
			t.Errorf("Expected Tech 1 to be available (New Format Skill)")
		}
		if !availableIDs[3] {
			t.Errorf("Expected Tech 3 to be available (Old Format Skill)")
		}
		if availableIDs[2] {
			t.Errorf("Expected Tech 2 to be unavailable (Wrong Skill)")
		}
		if availableIDs[4] {
			t.Errorf("Expected Tech 4 to be unavailable (No Skill)")
		}

		// Check Unavailable Reasons
		for _, tech := range response.Data.Unavailable {
			if tech.ID == 2 {
				if tech.Reason != "skill_mismatch" {
					t.Errorf("Expected Tech 2 reason to be skill_mismatch, got %s", tech.Reason)
				}
			}
		}
	})
}
