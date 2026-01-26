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

		// Should be 2 records (01-01 and 01-31)
		// 注意：我们在 BatchSetSchedule 测试中也插入了一些数据，
		// 但是 TestBatchSetSchedule 和 TestGetSchedules 是两个不同的测试函数，
		// 它们应该使用不同的 DB 实例或者清理数据。
		// 这里的问题在于 setupScheduleTestDB 使用了 shared cache 的内存数据库 "file::memory:?cache=shared"
		// 这意味着如果两个测试并行运行或者没有正确重置，数据会共享。
		// 在 TestBatchSetSchedule 中我们插入了:
		// 1. tech_id=1, date=2026-01-31 (Create New Schedule)
		// 2. tech_id=2, date=2026-02-01
		// 3. tech_id=3, date=2026-02-02
		// 4. tech_id=4,5 date=03-01,03-02
		//
		// 在 TestGetSchedules 中我们又插入了:
		// 1. tech_id=1, date=2026-01-01
		// 2. tech_id=1, date=2026-01-31 (重复插入?) -> 如果是同一个DB实例，可能会有冲突或者叠加
		//
		// 观察失败日志: Dates: [2026-01-31 00:00:00 2026-01-01 00:00:00 2026-01-31 00:00:00]
		// 我们有两条 2026-01-31 的记录！
		// 一条来自 TestBatchSetSchedule (Case 1)，一条来自 TestGetSchedules (Pre-create data 2)
		//
		// 修复方案：给 TestGetSchedules 使用一个唯一的数据库连接名，或者在 setup 中先清空表。
		// 简单起见，我们在 TestGetSchedules 中使用不一样的内存数据库名。

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
