package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupScheduleTestDB() *gorm.DB {
	// Use in-memory SQLite for testing
	// Use unique memory name to avoid conflicts if needed, but shared cache is fine
	d, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	d.AutoMigrate(&models.Schedule{})
	return d
}

func TestBatchSetSchedule(t *testing.T) {
	// Setup DB
	testDB := setupScheduleTestDB()
	originalDB := db.DB
	db.DB = testDB // Replace global DB
	defer func() { db.DB = originalDB }()

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/schedules/batch", BatchSetSchedule)

	// Test Case 1: Create new schedule with is_available = false
	t.Run("Create New Schedule (Leave)", func(t *testing.T) {
		reqBody := BatchSetScheduleRequest{
			TechIDs:     []uint{1},
			Dates:       []string{"2026-01-31"},
			IsAvailable: false,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/schedules/batch", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var schedule models.Schedule
		checkDate, _ := time.Parse("2006-01-02", "2026-01-31")
		err := testDB.Where("tech_id = ? AND date = ?", 1, datatypes.Date(checkDate)).First(&schedule).Error
		if err != nil {
			t.Errorf("Failed to find schedule: %v", err)
		}
		if schedule.IsAvailable != false {
			t.Errorf("Expected IsAvailable to be false, got true")
		}
	})

	// Test Case 2: Update existing schedule from true to false
	t.Run("Update Existing Schedule to Leave", func(t *testing.T) {
		// Pre-create a schedule
		date, _ := time.Parse("2006-01-02", "2026-02-01")
		testDB.Create(&models.Schedule{
			TechID:      2,
			Date:        datatypes.Date(date),
			IsAvailable: true,
		})

		reqBody := BatchSetScheduleRequest{
			TechIDs:     []uint{2},
			Dates:       []string{"2026-02-01"},
			IsAvailable: false,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/schedules/batch", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var schedule models.Schedule
		testDB.Where("tech_id = ? AND date = ?", 2, datatypes.Date(date)).First(&schedule)
		if schedule.IsAvailable != false {
			t.Errorf("Expected IsAvailable to be updated to false, got true")
		}
	})

	// Test Case 3: Update existing schedule from false to true
	t.Run("Update Existing Schedule to Work", func(t *testing.T) {
		// Pre-create a schedule
		date, _ := time.Parse("2006-01-02", "2026-02-02")
		testDB.Create(&models.Schedule{
			TechID:      3,
			Date:        datatypes.Date(date),
			IsAvailable: false,
		})

		reqBody := BatchSetScheduleRequest{
			TechIDs:     []uint{3},
			Dates:       []string{"2026-02-02"},
			IsAvailable: true,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/schedules/batch", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var schedule models.Schedule
		testDB.Where("tech_id = ? AND date = ?", 3, datatypes.Date(date)).First(&schedule)
		if schedule.IsAvailable != true {
			t.Errorf("Expected IsAvailable to be updated to true, got false")
		}
	})

	// Test Case 4: Multiple Techs and Dates
	t.Run("Batch Multiple Techs and Dates", func(t *testing.T) {
		reqBody := BatchSetScheduleRequest{
			TechIDs:     []uint{4, 5},
			Dates:       []string{"2026-03-01", "2026-03-02"},
			IsAvailable: false,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/schedules/batch", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var count int64
		d1, _ := time.Parse("2006-01-02", "2026-03-01")
		d2, _ := time.Parse("2006-01-02", "2026-03-02")
		testDB.Model(&models.Schedule{}).Where("tech_id IN ? AND date IN ? AND is_available = ?", []uint{4, 5}, []time.Time{d1, d2}, false).Count(&count)
		if count != 4 {
			t.Errorf("Expected 4 records, got %d", count)
		}
	})
}

func TestGetSchedules(t *testing.T) {
	// Setup DB with a unique name to avoid sharing data with TestBatchSetSchedule
	// Use ?cache=shared&mode=memory is enough for uniqueness per process run if not parallel,
	// but to be safe with shared cache, we can use a different path
	d, err := gorm.Open(sqlite.Open("file:TestGetSchedules?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	d.AutoMigrate(&models.Schedule{})
	testDB := d

	originalDB := db.DB
	db.DB = testDB // Replace global DB
	defer func() { db.DB = originalDB }()

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/api/schedules", GetSchedules)

	// Pre-create data
	// 1. Data on start_date
	d1, _ := time.Parse("2006-01-02", "2026-01-01")
	// Use map to ensure is_available=false is saved
	testDB.Model(&models.Schedule{}).Create(map[string]interface{}{
		"tech_id":      1,
		"date":         d1,
		"is_available": false,
	})

	// 2. Data on end_date
	d2, _ := time.Parse("2006-01-02", "2026-01-31")
	testDB.Model(&models.Schedule{}).Create(map[string]interface{}{
		"tech_id":      1,
		"date":         d2,
		"is_available": false,
	})

	// 3. Data out of range
	d3, _ := time.Parse("2006-01-02", "2026-02-02") // Change to 02-02 to be safely out of range
	testDB.Model(&models.Schedule{}).Create(map[string]interface{}{
		"tech_id":      1,
		"date":         d3,
		"is_available": false,
	})

	// Test Case 1: Query range [2026-01-01, 2026-01-31]
	t.Run("Query Range Inclusive", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Simulate frontend request, note end_date is 01-31
		req, _ := http.NewRequest("GET", "/api/schedules?start_date=2026-01-01&end_date=2026-01-31", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response struct {
			Code int               `json:"code"`
			Data []models.Schedule `json:"data"`
			Msg  string            `json:"msg"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		if len(response.Data) != 2 {
			// For debugging, print actual returned dates
			var dates []string
			for _, s := range response.Data {
				dates = append(dates, time.Time(s.Date).Format("2006-01-02"))
			}
			t.Errorf("Expected 2 records, got %d. Dates: %v", len(response.Data), dates)
		}

		foundEnd := false
		for _, s := range response.Data {
			if time.Time(s.Date).Format("2006-01-02") == "2026-01-31" {
				foundEnd = true
				if s.IsAvailable != false {
					t.Errorf("Expected 2026-01-31 to be unavailable (false), got true")
				}
			}
		}
		if !foundEnd {
			t.Errorf("Expected to find record for 2026-01-31")
		}
	})
}

func TestScheduleIntegration(t *testing.T) {
	// Setup DB
	d, err := gorm.Open(sqlite.Open("file:TestScheduleIntegration?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	d.AutoMigrate(&models.Schedule{}, &models.Technician{})
	testDB := d
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/schedules/batch", BatchSetSchedule)
	r.GET("/api/schedules", GetSchedules)

	// Create Technicians
	techIDs := []uint{2, 3, 4, 5, 6, 7, 8}
	for _, id := range techIDs {
		testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: id}, Name: "Tech " + string(rune(id+65))}) // A, B, ...
	}

	// Step 1: Set 01-31 Available (Work)
	t.Run("Step 1: Set 01-31 Available", func(t *testing.T) {
		reqBody := BatchSetScheduleRequest{
			TechIDs:     techIDs,
			Dates:       []string{"2026-01-31"},
			IsAvailable: true,
		}
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/schedules/batch", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	// Step 2: Verify 01-31 Available
	t.Run("Step 2: Verify 01-31 Available", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Assuming default pagination or enough limit, or filter by date
		req, _ := http.NewRequest("GET", "/api/schedules?start_date=2026-01-01&end_date=2026-01-31", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response struct {
			Code int               `json:"code"`
			Data []models.Schedule `json:"data"`
			Msg  string            `json:"msg"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal: %v", err)
		}

		// Check if we have records for 2026-01-31 and they are available
		count := 0
		for _, s := range response.Data {
			dateStr := time.Time(s.Date).Format("2006-01-02")
			if dateStr == "2026-01-31" {
				count++
				if !s.IsAvailable {
					t.Errorf("Expected tech %d to be available on 2026-01-31", s.TechID)
				}
			}
		}
		if count != len(techIDs) {
			t.Errorf("Expected %d records for 2026-01-31, got %d", len(techIDs), count)
		}
	})

	// Step 3: Set 01-31 Unavailable (Leave)
	t.Run("Step 3: Set 01-31 Unavailable", func(t *testing.T) {
		reqBody := BatchSetScheduleRequest{
			TechIDs:     techIDs,
			Dates:       []string{"2026-01-31"},
			IsAvailable: false,
		}
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/schedules/batch", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	// Step 4: Verify 01-31 Unavailable
	t.Run("Step 4: Verify 01-31 Unavailable", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/schedules?start_date=2026-01-01&end_date=2026-01-31", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response struct {
			Code int               `json:"code"`
			Data []models.Schedule `json:"data"`
			Msg  string            `json:"msg"`
		}
		json.Unmarshal(w.Body.Bytes(), &response)

		count := 0
		for _, s := range response.Data {
			dateStr := time.Time(s.Date).Format("2006-01-02")
			if dateStr == "2026-01-31" {
				count++
				// The requirement says "1. 设置01-31休假，测试01-31为工作"
				// This implies user might WANT to see if it fails to update, OR expects it to be Working?
				// Logic dictates: Set Leave -> Expect Leave.
				// If I strictly follow "test 01-31 is work", then I should assert IsAvailable == true.
				// But that would mean the Set Leave failed.
				// I will assume "test 01-31 is work" is a typo for "test 01-31 is leave" OR "test 01-31 status".
				// I will assert it is Unavailable (false).
				if s.IsAvailable {
					t.Errorf("Expected tech %d to be unavailable on 2026-01-31", s.TechID)
				}
			}
		}
		if count != len(techIDs) {
			t.Errorf("Expected %d records for 2026-01-31, got %d", len(techIDs), count)
		}
	})
}
