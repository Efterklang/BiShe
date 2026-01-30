package handlers

import (
	"net/http"
	"time"

	"server/internal/db"
	"server/internal/models"
	"server/internal/repo"
	"server/pkg/config"

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

	sDate, _ := time.Parse("2006-01-02", startDateStr)
	eDate, _ := time.Parse("2006-01-02", endDateStr)

	query := db.DB.Model(&models.Schedule{}).
		Preload("Technician").
		Where("date >= ? AND date <= ?", datatypes.Date(sDate), datatypes.Date(eDate))

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

	// 1. 预处理日期：只解析一次，避免在循环中重复解析
	// 同时在这里做格式校验，如果有一个日期格式不对，直接报错，无需开启事务
	var targetDates []datatypes.Date
	for _, dateStr := range req.Dates {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期格式错误: " + dateStr})
			return
		}
		targetDates = append(targetDates, datatypes.Date(parsedDate))
	}

	schedules := make([]map[string]interface{}, 0, len(req.TechIDs)*len(req.Dates))

	for _, techID := range req.TechIDs {
		for _, date := range targetDates {
			schedules = append(schedules, map[string]interface{}{
				"tech_id":      techID,
				"date":         date,
				"is_available": req.IsAvailable,
			})
		}
	}

	// 批量Upsert (SQLite)
	// 注意: 此操作要求 schedules 表在 (tech_id, date) 上有唯一索引
	// 使用 map 插入以避免 GORM 对 bool 零值(false)应用 default tag
	err := repo.Schedule.BatchUpsertSchedules(schedules)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "批量设置排班失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "排班设置成功"})
}

// GetAvailableTechnicians 获取指定时间段的可用技师列表
// Query Params:
//   - start_time (RFC3339格式，必填)
//   - service_id (服务项目ID，必填)
//
// 筛选逻辑:
//  1. 检查技师是否具备该服务项目所需的技能
//  2. 检查技师当天是否在岗（排班状态）
//  3. 检查技师在该时间段是否有冲突预约
//
// Response:
//   - available: 具备技能且时间空闲的技师
//   - service: 请求的服务项目信息
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

	// 1. 获取服务项目信息，计算结束时间
	var service models.ServiceProduct
	if err := db.DB.First(&service, serviceIDStr).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Service not found"})
		return
	}

	endTime := startTime.Add(time.Duration(service.Duration) * time.Minute)
	// 使用 UTC 格式化日期，确保与 datatypes.Date 存储一致
	dateStr := startTime.UTC().Format("2006-01-02")

	// 2. 获取具有服务技能的技师
	skilledTechnicians, err := repo.Technician.GetTechniciansWithSkill(service.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取技师失败", "error": err.Error()})
		return
	}

	// 3. 获取空闲技师
	freeTechnicians, err := repo.Schedule.GetAvailableTechs(dateStr, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取排班失败", "error": err.Error()})
		return
	}

	// 构建空闲技师Map加速查找
	freeTechMap := make(map[uint]bool)
	for _, t := range freeTechnicians {
		freeTechMap[t.ID] = true
	}

	// 获取当天请假/不可用的技师ID，用于区分Reason
	leaveTechMap, err := repo.Schedule.GetUnavailableTechIDs(dateStr)
	if err != nil {
		// 容错：如果获取失败，默认为空Map
		leaveTechMap = make(map[uint]bool)
	}

	var availableTechnicians []models.Technician
	var unavailableTechnicians []models.Technician

	// 4. 分类技师
	for i := range skilledTechnicians {
		tech := &skilledTechnicians[i] // 使用指针

		if freeTechMap[tech.ID] {
			availableTechnicians = append(availableTechnicians, *tech)
		} else {
			// 确定不可用原因
			if leaveTechMap[tech.ID] {
				tech.Reason = "leave"
			} else {
				tech.Reason = "busy"
			}
			unavailableTechnicians = append(unavailableTechnicians, *tech)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"available":   availableTechnicians,
			"unavailable": unavailableTechnicians,
			"service": gin.H{
				"id":       service.ID,
				"name":     service.Name,
				"price":    service.Price,
				"duration": service.Duration,
			},
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
	schedule, err := repo.Schedule.GetByTechAndDate(techIDStr, dateStr)

	// 如果没有排班记录，默认为可用
	isAvailable := true
	if err == nil {
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

// GetTimeSlotsAvailability 获取某日所有时间段的可用性
// Query Params: date (YYYY-MM-DD), service_id
func GetTimeSlotsAvailability(c *gin.Context) {
	dateStr := c.Query("date")
	serviceIDStr := c.Query("service_id")

	if dateStr == "" || serviceIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "date and service_id are required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Invalid date format"})
		return
	}

	var service models.ServiceProduct
	if err := db.DB.First(&service, serviceIDStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Service not found"})
		return
	}

	// 1. 预处理：筛选具备技能的技师ID
	skilledTechs, err := repo.Technician.GetTechniciansWithSkill(service.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Failed to fetch skilled technicians"})
		return
	}
	var skilledTechIDs []uint
	for _, tech := range skilledTechs {
		skilledTechIDs = append(skilledTechIDs, tech.ID)
	}

	// 2. 预处理：获取当天排班（请假情况）
	// datatypes.Date 存储为 "YYYY-MM-DD"，查询时直接使用 dateStr 即可
	// checkDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	unavailableScheduleTechs, _ := repo.Schedule.GetUnavailableTechIDs(dateStr)

	// 3. 预处理：获取当天所有预约
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	var appointments []models.Appointment
	db.DB.Where("start_time >= ? AND start_time < ? AND status = ?", startOfDay, endOfDay, "pending").
		Find(&appointments)

	// 生成时间槽 (从配置读取)
	openTime := time.Date(date.Year(), date.Month(), date.Day(), config.GlobalBusinessHours.OpenHour, config.GlobalBusinessHours.OpenMinute, 0, 0, config.GlobalBusinessHours.TimeLocation)
	closeTime := time.Date(date.Year(), date.Month(), date.Day(), config.GlobalBusinessHours.CloseHour, config.GlobalBusinessHours.CloseMinute, 0, 0, config.GlobalBusinessHours.TimeLocation)

	type TimeSlot struct {
		Time           string `json:"time"`            // "10:00"
		Status         string `json:"status"`          // "available", "waitlist", "closed"
		AvailableCount int    `json:"available_count"` // 可用技师数
		StartTime      string `json:"start_time"`      // ISO String
	}

	var slots []TimeSlot

	for t := openTime; t.Before(closeTime); t = t.Add(config.GlobalBusinessHours.SlotInterval) {
		slotEndTime := t.Add(time.Duration(service.Duration) * time.Minute)

		// 如果服务结束时间超过打烊时间，则不可约
		if slotEndTime.After(closeTime) {
			slots = append(slots, TimeSlot{
				Time:           t.Format("15:04"),
				Status:         "closed",
				AvailableCount: 0,
				StartTime:      t.Format(time.RFC3339),
			})
			continue
		}

		// 检查过去时间
		if t.Before(time.Now()) {
			slots = append(slots, TimeSlot{
				Time:           t.Format("15:04"),
				Status:         "closed", // Past time
				AvailableCount: 0,
				StartTime:      t.Format(time.RFC3339),
			})
			continue
		}

		// 计算该槽位可用技师
		availableCount := 0
		for _, techID := range skilledTechIDs {
			// 1. 检查排班
			if unavailableScheduleTechs[techID] {
				continue
			}

			// 2. 检查预约冲突
			isBusy := false
			for _, appt := range appointments {
				if appt.TechID == techID {
					// 检查时间重叠
					// max(start1, start2) < min(end1, end2)
					if t.Before(appt.EndTime) && slotEndTime.After(appt.StartTime) {
						isBusy = true
						break
					}
				}
			}
			if isBusy {
				continue
			}

			availableCount++
		}

		status := "waitlist"
		if availableCount > 0 {
			status = "available"
		}

		slots = append(slots, TimeSlot{
			Time:           t.Format("15:04"),
			Status:         status,
			AvailableCount: availableCount,
			StartTime:      t.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": slots,
		"msg":  "success",
	})
}
