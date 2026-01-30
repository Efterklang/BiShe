package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"server/internal/db"
	"server/internal/models"
)

// setupTestDB 初始化一个用于测试的内存数据库
func setupTestDB() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// 自动迁移模型
	err = database.AutoMigrate(
		&models.Technician{},
		&models.ServiceProduct{},
		&models.Appointment{},
	)
	if err != nil {
		return nil, err
	}
	return database, nil
}

// TestListTechnicians_SkillNamesMapping 测试 ListTechnicians 接口的 skill_names 映射功能
func TestListTechnicians_SkillNamesMapping(t *testing.T) {
	// 1. 初始化测试数据库
	testDB, err := setupTestDB()
	if err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}

	// 保存原始DB并在测试结束后恢复
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// 2. 准备测试数据
	// 创建服务产品
	service1 := models.ServiceProduct{Name: "精油SPA", Duration: 90, Price: 298.0, IsActive: true}
	service2 := models.ServiceProduct{Name: "中式推拿", Duration: 60, Price: 168.0, IsActive: true}
	testDB.Create(&service1)
	testDB.Create(&service2)

	// 创建技师，关联服务ID
	skills1 := []uint{service1.ID, service2.ID}
	skillsData1, _ := json.Marshal(skills1)
	tech1 := models.Technician{Name: "张师傅", Skills: skillsData1, Status: 0, AverageRating: 4.5}

	// 创建没有技能的技师
	var skillsData2 []byte // 保持为 nil
	tech2 := models.Technician{Name: "李师傅", Skills: skillsData2, Status: 0, AverageRating: 4.0}

	testDB.Create(&tech1)
	testDB.Create(&tech2)

	// 3. 初始化 Gin 路由器并注册路由
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 使用闭包将 testDB 注入到处理器中
	router.GET("/api/technicians", func(c *gin.Context) {
		ListTechnicians(c) // 将测试数据库实例传入处理器
	})

	// 4. 创建测试请求
	req, _ := http.NewRequest("GET", "/api/technicians", nil)
	w := httptest.NewRecorder()

	// 5. 执行请求
	router.ServeHTTP(w, req)

	// 6. 验证响应
	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 %d, 实际状态码 %d", http.StatusOK, w.Code)
	}

	var response struct {
		Code int `json:"code"`
		Data []struct {
			ID            uint     `json:"id"`
			Name          string   `json:"name"`
			Status        int      `json:"status"`
			AverageRating float32  `json:"average_rating"`
			SkillNames    []string `json:"skill_names"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应JSON失败: %v", err)
	}

	// 7. 验证响应数据
	if len(response.Data) != 2 {
		t.Fatalf("期望返回 %d 个技师, 实际返回 %d 个", 2, len(response.Data))
	}

	// 验证第一个技师（有技能）
	if response.Data[0].Name != "张师傅" {
		t.Errorf("第一个技师名称错误: 期望 '张师傅', 实际 '%s'", response.Data[0].Name)
	}
	if len(response.Data[0].SkillNames) != 2 {
		t.Errorf("第一个技师技能数量错误: 期望 2, 实际 %d", len(response.Data[0].SkillNames))
	}
	// 验证技能名称是否正确映射
	expectedSkills := []string{"精油SPA", "中式推拿"}
	if !slicesEqual(response.Data[0].SkillNames, expectedSkills) {
		t.Errorf("第一个技师技能名称错误: 期望 %v, 实际 %v", expectedSkills, response.Data[0].SkillNames)
	}

	// 验证第二个技师（无技能）
	if response.Data[1].Name != "李师傅" {
		t.Errorf("第二个技师名称错误: 期望 '李师傅', 实际 '%s'", response.Data[1].Name)
	}
	if response.Data[1].SkillNames == nil {
		t.Errorf("第二个技师的 skill_names 不应该是 null")
	}
	if len(response.Data[1].SkillNames) != 0 {
		t.Errorf("第二个技师技能数量错误: 期望 0, 实际 %d", len(response.Data[1].SkillNames))
	}
}

// slicesEqual 比较两个字符串切片是否相等
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
