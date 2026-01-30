package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"server/internal/db"
	"server/internal/models"

	// "server/internal/repo" // Implicitly used by handlers

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

	// 创建技术人员及其技能（仅支持新格式：服务ID数组）
	// 技师1: 拥有按摩技能 (新格式: [1])
	skills1, _ := json.Marshal([]uint{1})
	testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: 1}, Name: "技师1", Skills: skills1})

	// 技师2: 拥有面部护理技能 (新格式: [2])
	skills2, _ := json.Marshal([]uint{2})
	testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: 2}, Name: "技师2", Skills: skills2})

	// 技师3: 拥有按摩技能 (新格式: [1])
	skills3, _ := json.Marshal([]uint{1})
	testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: 3}, Name: "技师3", Skills: skills3})

	// 技师4: 无技能
	testDB.Create(&models.Technician{BaseModel: models.BaseModel{ID: 4}, Name: "技师4"})

	// 测试用例：筛选按摩服务（服务ID 1）
	t.Run("按按摩技能筛选", func(t *testing.T) {
		w := httptest.NewRecorder()
		startTime := time.Now().Add(time.Hour).Format(time.RFC3339)

		params := url.Values{}
		params.Add("service_id", "1")
		params.Add("start_time", startTime)

		req, _ := http.NewRequest("GET", "/api/schedules/available-technicians?"+params.Encode(), nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码200，实际状态码: %d. 响应体: %s", w.Code, w.Body.String())
		}

		var response struct {
			Data struct {
				Available   []models.Technician `json:"available"`
				Unavailable []models.Technician `json:"unavailable"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("解析响应失败: %v", err)
		}

		// 期望可预约：技师1、技师3（都拥有按摩技能）
		availableIDs := make(map[uint]bool)
		for _, t := range response.Data.Available {
			availableIDs[t.ID] = true
		}

		if !availableIDs[1] {
			t.Errorf("期望技师1可预约（拥有按摩技能）")
		}
		if !availableIDs[3] {
			t.Errorf("期望技师3可预约（拥有按摩技能）")
		}
		if availableIDs[2] {
			t.Errorf("期望技师2不可预约（技能不匹配）")
		}
		if availableIDs[4] {
			t.Errorf("期望技师4不可预约（无技能）")
		}

		// 检查不可预约原因
		for _, tech := range response.Data.Unavailable {
			if tech.ID == 2 {
				if tech.Reason != "skill_mismatch" {
					t.Errorf("期望技师2的原因为技能不匹配，实际原因: %s", tech.Reason)
				}
			}
		}
	})
}
