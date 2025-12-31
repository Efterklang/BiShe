package handlers

import (
	"net/http"
	"time"

	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// GetSchedules 获取排班列表 (支持月视图和日视图的数据需求)
// Query Params: start_date (YYYY-MM-DD), end_date (YYYY-MM-DD), tech_ids (1,2,3)
func GetSchedules(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// 默认查询未来 30 天
	if startDateStr == "" {
		startDateStr = time.Now().Format("2006-01-02")
	}
	if endDateStr == "" {
		endDateStr = time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	}

	var schedules []models.Schedule
	query := db.DB.Model(&models.Schedule{}).
		Preload("Technician").
		Where("date >= ? AND date <= ?", startDateStr, endDateStr)

	// 支持按技师筛选
	techIDs := c.QueryArray("tech_ids[]") // 前端传 tech_ids[]=1&tech_ids[]=2
	if len(techIDs) > 0 {
		query = query.Where("tech_id IN ?", techIDs)
	}

	if err := query.Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取排班失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": schedules,
		"msg":  "success",
	})
}

// BatchSetScheduleRequest 批量设置排班请求体
type BatchSetScheduleRequest struct {
	TechIDs     []uint   `json:"tech_ids" binding:"required"`
	Dates       []string `json:"dates" binding:"required"` // ["2023-10-01", "2023-10-02"]
	IsAvailable bool     `json:"is_available"`             // true: 上班, false: 请假
}

// BatchSetSchedule 批量设置技师排班 (用于请假或安排工作)
func BatchSetSchedule(c *gin.Context) {
	var req BatchSetScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	// 开启事务
	tx := db.DB.Begin()

	for _, techID := range req.TechIDs {
		for _, dateStr := range req.Dates {
			// 解析日期
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期格式错误: " + dateStr})
				return
			}

			var schedule models.Schedule
			// 查找是否存在记录，存在则更新，不存在则创建 (Upsert)
			result := tx.Where("tech_id = ? AND date = ?", techID, date).First(&schedule)

			if result.RowsAffected == 0 {
				// 创建新记录
				schedule = models.Schedule{
					TechID:      techID,
					Date:        date,
					IsAvailable: req.IsAvailable,
					TimeSlots:   datatypes.JSON([]byte("[]")), // 初始化为空数组
				}
				if err := tx.Create(&schedule).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建排班失败"})
					return
				}
			} else {
				// 更新现有记录状态
				if err := tx.Model(&schedule).Update("is_available", req.IsAvailable).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新排班失败"})
					return
				}
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "排班设置成功"})
}
