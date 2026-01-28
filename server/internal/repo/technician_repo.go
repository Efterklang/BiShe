package repo

import (
	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"
)

type TechnicianRepository interface {
	GetTechniciansWithSkill(serviceID uint) ([]models.Technician, error)
}

type TechnicianRepo struct{}

var Technician TechnicianRepository = &TechnicianRepo{}

// GetTechniciansWithSkill 筛选具备指定技能（serviceID）的技师
func (r *TechnicianRepo) GetTechniciansWithSkill(serviceID uint) ([]models.Technician, error) {
	// 定义返回结果集
	var skilledTechs []models.Technician

	// 核心SQL逻辑：
	// 1. json_valid(Skills)：确保Skills是合法JSON格式
	// 2. json_type(Skills) = 'array'：确保Skills是JSON数组（而非单个值/对象）
	// 3. EXISTS子查询：遍历JSON数组，判断是否包含目标serviceID
	// 4. CAST(value AS INTEGER)：兼容偶发的字符串格式数值（如["1","2"]），保证匹配精准
	query := `
		SELECT * FROM technicians
		WHERE json_valid(Skills) = 1
		  AND json_type(Skills) = 'array'
		  AND EXISTS (
			SELECT 1 FROM json_each(Skills)
			WHERE CAST(value AS INTEGER) = ?
		  )
	`

	// 执行原生SQL查询，将serviceID作为参数传入（防止SQL注入）
	// Raw执行原生SQL，Scan将结果映射到models.Technician结构体切片
	if err := db.DB.Raw(query, serviceID).Scan(&skilledTechs).Error; err != nil {
		return nil, err
	}

	return skilledTechs, nil
}
