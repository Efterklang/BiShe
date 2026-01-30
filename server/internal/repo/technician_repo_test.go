package repo

import (
	"encoding/json"
	"testing"

	"server/internal/db"
	"server/internal/models"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// 使用内存数据库
	d, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	// 迁移 Technician 模型
	err = d.AutoMigrate(&models.Technician{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	return d
}

func TestGetTechniciansWithSkill(t *testing.T) {
	// 1. 设置测试数据库
	testDB := setupTestDB(t)

	// 替换全局 DB 实例
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// 2. 准备测试数据
	targetSkillID := uint(101)
	otherSkillID := uint(202)

	// Case 1: 具备技能的技师 (应该包含)
	skills1, _ := json.Marshal([]uint{targetSkillID, 303})
	techWithSkill := models.Technician{
		Name:   "Tech With Skill",
		Skills: datatypes.JSON(skills1),
	}
	if err := testDB.Create(&techWithSkill).Error; err != nil {
		t.Fatalf("Failed to create techWithSkill: %v", err)
	}

	// Case 2: 不具备技能的技师 (应该排除)
	skills2, _ := json.Marshal([]uint{otherSkillID, 303})
	techWithoutSkill := models.Technician{
		Name:   "Tech Without Skill",
		Skills: datatypes.JSON(skills2),
	}
	if err := testDB.Create(&techWithoutSkill).Error; err != nil {
		t.Fatalf("Failed to create techWithoutSkill: %v", err)
	}

	// Case 3: 已删除的具备技能的技师 (应该排除)
	skills3, _ := json.Marshal([]uint{targetSkillID})
	techDeleted := models.Technician{
		Name:   "Tech Deleted",
		Skills: datatypes.JSON(skills3),
	}
	if err := testDB.Create(&techDeleted).Error; err != nil {
		t.Fatalf("Failed to create techDeleted: %v", err)
	}
	// 执行软删除
	if err := testDB.Delete(&techDeleted).Error; err != nil {
		t.Fatalf("Failed to delete techDeleted: %v", err)
	}

	// Case 4: 具备技能的技师，但存储为字符串格式 (应该包含，因为 SQL 中有 CAST)
	// 例如 ["101", "404"]
	skills4, _ := json.Marshal([]string{"101", "404"})
	techStringSkill := models.Technician{
		Name:   "Tech String Skill",
		Skills: datatypes.JSON(skills4),
	}
	if err := testDB.Create(&techStringSkill).Error; err != nil {
		t.Fatalf("Failed to create techStringSkill: %v", err)
	}

	// 3. 执行测试
	repo := &TechnicianRepo{}
	results, err := repo.GetTechniciansWithSkill(targetSkillID)

	// 4. 验证结果
	if err != nil {
		t.Fatalf("GetTechniciansWithSkill returned error: %v", err)
	}

	// 预期找到 2 个技师 (techWithSkill 和 techStringSkill)
	expectedCount := 2
	if len(results) != expectedCount {
		t.Errorf("Expected %d technicians, got %d", expectedCount, len(results))
	}

	// 验证 ID 是否匹配
	foundMap := make(map[uint]bool)
	for _, tech := range results {
		foundMap[tech.ID] = true
	}

	if !foundMap[techWithSkill.ID] {
		t.Errorf("TechWithSkill (ID: %d) not found", techWithSkill.ID)
	}
	if !foundMap[techStringSkill.ID] {
		t.Errorf("TechStringSkill (ID: %d) not found", techStringSkill.ID)
	}
	if foundMap[techWithoutSkill.ID] {
		t.Errorf("TechWithoutSkill (ID: %d) should NOT be found", techWithoutSkill.ID)
	}
	if foundMap[techDeleted.ID] {
		t.Errorf("TechDeleted (ID: %d) should NOT be found", techDeleted.ID)
	}
}
