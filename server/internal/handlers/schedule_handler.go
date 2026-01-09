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

// GetAvailableTechnicians 获取指定时间段的可用技师列表
// Query Params: start_time (RFC3339), service_id (必填，用于计算结束时间)
func GetAvailableTechnicians(c *gin.Context) {
	startTimeStr := c.Query("start_time")
	serviceIDStr := c.Query("service_id")

	if startTimeStr == "" || serviceIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "start_time and service_id are required"})
		return
	}

	// 解析开始时间
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Invalid start_time format"})
		return
	}

	// 获取服务项目信息，计算结束时间
	var service models.ServiceItem
	if err := db.DB.First(&service, serviceIDStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Service not found"})
		return
	}

	endTime := startTime.Add(time.Duration(service.Duration) * time.Minute)

	// 获取所有技师
	var allTechnicians []models.Technician
	if err := db.DB.Find(&allTechnicians).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Failed to fetch technicians"})
		return
	}

	// 检查排班状态（是否在岗）
	checkDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	var schedules []models.Schedule
	db.DB.Where("date = ?", checkDate).Find(&schedules)

	scheduleMap := make(map[uint]bool)
	for _, s := range schedules {
		scheduleMap[s.TechID] = s.IsAvailable
	}

	// 查询该时间段有冲突的技师（有 pending 预约）
	var conflictAppointments []models.Appointment
	db.DB.Where("status = ? AND start_time < ? AND end_time > ?", "pending", endTime, startTime).
		Find(&conflictAppointments)

	busyTechMap := make(map[uint]bool)
	for _, appt := range conflictAppointments {
		busyTechMap[appt.TechID] = true
	}

	// 筛选可用技师
	var availableTechnicians []models.Technician
	var unavailableTechnicians []models.Technician

	for _, tech := range allTechnicians {
		// 检查是否在岗（默认在岗）
		isAvailable, exists := scheduleMap[tech.ID]
		if exists && !isAvailable {
			// 休息中
			unavailableTechnicians = append(unavailableTechnicians, tech)
			continue
		}

		// 检查是否有冲突预约
		if busyTechMap[tech.ID] {
			// 忙碌中
			unavailableTechnicians = append(unavailableTechnicians, tech)
			continue
		}

		// 可用
		availableTechnicians = append(availableTechnicians, tech)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"available":   availableTechnicians,
			"unavailable": unavailableTechnicians,
			"total":       len(allTechnicians),
		},
		"msg": "success",
	})
}

// GetTechnicianScheduleDetail 获取技师排班详情（包含预约信息）
// Query Params: tech_id (required), date (YYYY-MM-DD, required)
func GetTechnicianScheduleDetail(c *gin.Context) {
	techIDStr := c.Query("tech_id")
	dateStr := c.Query("date")

	if techIDStr == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "tech_id and date are required"})
		return
	}

	// 解析日期
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Invalid date format"})
		return
	}

	// 获取排班信息
	var schedule models.Schedule
	scheduleErr := db.DB.Where("tech_id = ? AND date = ?", techIDStr, date).First(&schedule).Error

	// 如果没有排班记录，默认为可用
	isAvailable := true
	if scheduleErr == nil {
		isAvailable = schedule.IsAvailable
	}

	// 获取当天该技师的预约列表
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var appointments []models.Appointment
	if err := db.DB.Where("tech_id = ? AND start_time >= ? AND start_time < ?", techIDStr, startOfDay, endOfDay).
		Preload("Member").
		Preload("ServiceItem").
		Order("start_time ASC").
		Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Failed to get appointments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"date":         dateStr,
			"tech_id":      techIDStr,
			"is_available": isAvailable,
			"appointments": appointments,
		},
		"msg": "success",
	})
}
