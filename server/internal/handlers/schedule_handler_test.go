package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"server/internal/db"
	"server/internal/models"
	"server/internal/repo"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockTechnicianRepo to override GetTechniciansWithSkill
type MockTechnicianRepo struct {
	SkilledTechs []models.Technician
}

func (m *MockTechnicianRepo) GetTechniciansWithSkill(serviceID uint) ([]models.Technician, error) {
	return m.SkilledTechs, nil
}

func TestGetAvailableTechnicians(t *testing.T) {
	// 1. 设置内存数据库
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法打开测试数据库: %v", err)
	}

	// 迁移所有相关模型
	err = testDB.AutoMigrate(
		&models.Technician{},
		&models.ServiceProduct{},
		&models.Schedule{},
		&models.Appointment{},
		&models.Member{},
	)
	if err != nil {
		t.Fatalf("迁移失败: %v", err)
	}

	// 替换全局数据库
	originalDB := db.DB
	db.DB = testDB

	// Mock仓库
	originalTechRepo := repo.Technician
	mockRepo := &MockTechnicianRepo{}
	repo.Technician = mockRepo

	defer func() {
		db.DB = originalDB
		repo.Technician = originalTechRepo
	}()

	// 2. 准备数据
	// 服务项目
	service := models.ServiceProduct{Name: "深层按摩", Duration: 60, Price: 200}
	testDB.Create(&service)

	// 测试日期应该固定以避免某些数据库驱动程序中的时区/今天问题
	// 在测试设置中始终使用UTC
	now := time.Now().UTC()
	testDate := datatypes.Date(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC))

	// 技术人员
	// 技师1: 可预约 (有技能, 有排班, 无预约)
	tech1 := models.Technician{Name: "张师傅", Skills: datatypes.JSON([]byte(fmt.Sprintf("[%d]", service.ID)))}
	testDB.Create(&tech1)
	testDB.Create(&models.Schedule{TechID: tech1.ID, Date: testDate, IsAvailable: true})

	// 技师2: 忙碌 (有技能, 有排班, 有冲突预约)
	tech2 := models.Technician{Name: "李师傅", Skills: datatypes.JSON([]byte(fmt.Sprintf("[%d]", service.ID)))}
	testDB.Create(&tech2)
	testDB.Create(&models.Schedule{TechID: tech2.ID, Date: testDate, IsAvailable: true})

	startTime := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, time.UTC)
	testDB.Create(&models.Appointment{
		TechID:    tech2.ID,
		ServiceID: service.ID,
		MemberID:  1,
		StartTime: startTime.Add(-30 * time.Minute),
		EndTime:   startTime.Add(30 * time.Minute),
		Status:    "pending",
	})

	// 技师3: 请假 (有技能, 但在排班中请假)
	tech3 := models.Technician{Name: "王师傅", Skills: datatypes.JSON([]byte(fmt.Sprintf("[%d]", service.ID)))}
	testDB.Create(&tech3)
	// 使用原始SQL创建排班条目或确保datatypes.Date正是仓库所期望的
	testDB.Exec("INSERT INTO schedules (tech_id, date, is_available, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		tech3.ID, now.Format("2006-01-02"), false, now, now)

	// 设置模拟数据
	mockRepo.SkilledTechs = []models.Technician{tech1, tech2, tech3}

	// 技师4: 无技能 (不在模拟SkilledTechs中)
	tech4 := models.Technician{Name: "赵师傅", Skills: datatypes.JSON([]byte("[]"))}
	testDB.Create(&tech4)
	testDB.Create(&models.Schedule{TechID: tech4.ID, Date: testDate, IsAvailable: true})

	// 3. 设置路由器
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/technicians/available", GetAvailableTechnicians)

	// 4. 测试用例
	t.Run("成功用例", func(t *testing.T) {
		// 使用UTC查询以匹配数据库存储
		startTimeStr := startTime.UTC().Format(time.RFC3339)
		params := url.Values{}
		params.Add("start_time", startTimeStr)
		params.Add("service_id", fmt.Sprintf("%d", service.ID))

		targetURL := fmt.Sprintf("/technicians/available?%s", params.Encode())

		req, _ := http.NewRequest("GET", targetURL, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("期望状态码200，实际状态码: %d. 响应体: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Code int `json:"code"`
			Data struct {
				Available   []models.Technician `json:"available"`
				Unavailable []models.Technician `json:"unavailable"`
				Service     struct {
					ID   uint   `json:"id"`
					Name string `json:"name"`
				} `json:"service"`
			} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)

		// 断言可预约
		if len(resp.Data.Available) != 1 || resp.Data.Available[0].ID != tech1.ID {
			t.Errorf("期望技师1可预约，实际结果: %v", resp.Data.Available)
		}

		// 断言不可预约
		if len(resp.Data.Unavailable) != 2 {
			t.Errorf("期望2个不可预约技师(忙碌, 请假)，实际数量: %d", len(resp.Data.Unavailable))
		}

		// 验证原因
		foundBusy := false
		foundLeave := false
		for _, ut := range resp.Data.Unavailable {
			if ut.ID == tech2.ID && ut.Reason == "busy" {
				foundBusy = true
			}
			if ut.ID == tech3.ID && ut.Reason == "leave" {
				foundLeave = true
			}
		}
		if !foundBusy || !foundLeave {
			t.Errorf("原因映射失败. 找到忙碌: %v, 找到请假: %v", foundBusy, foundLeave)
		}
	})
}
