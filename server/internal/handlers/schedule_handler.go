package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	// 查询时使用 Preload 加载 Technician 信息
	// datatypes.Date 映射到数据库中的 DATE 类型（字符串），可以直接进行比较
	// 且不需要处理时区和时分秒问题，date >= '2026-01-01' AND date <= '2026-01-31' 可以正确包含边界
	// 注意：endDateStr 已经是 "YYYY-MM-DD" 格式，可以直接比较
	// 关键：GORM 的 datatypes.Date 实现了 Scanner/Valuer，在查询时会将 time.Time 或 string 转换为正确的格式
	// 但是在 Where 语句中，如果传入的是 string，数据库会进行字符串比较。
	// 对于 SQLite，DATE 类型本质上是 TEXT。
	// "2026-01-31" <= "2026-01-31" 是 true。
	// GORM 默认使用 time.Time 解析 date，所以 Where("date >= ? AND date <= ?", startDateStr, endDateStr) 会把 string 解析为 time.Time
	// 这样就变成了 date >= '2026-01-01 00:00:00+00:00' AND date <= '2026-01-31 00:00:00+00:00'
	// 而数据库里存的是 "2026-01-31"，字符串比较 '2026-01-31' <= '2026-01-31 00:00:00+00:00' 是 true
	// 等等，数据库里存的是什么？
	// 迁移脚本执行的是 `UPDATE schedules SET date = substr(date, 1, 10)`，所以数据库里存的是 "2026-01-31"
	// 测试数据使用 `datatypes.Date(time.Time)` 创建。GORM 在 SQLite 下 datatypes.Date 的 Valuer 实现：
	// return time.Time(date).Format("2006-01-02"), nil
	// 所以测试数据也是存为 "2026-01-31"。
	//
	// 查询参数 startDateStr 和 endDateStr 是字符串 "2026-01-01" 和 "2026-01-31"。
	// GORM Where 的行为：
	// 如果参数是 string，且没有手动 Cast，GORM 可能会直接传给数据库驱动。
	// 但如果模型字段是 time.Time，GORM 可能会尝试解析。这里模型字段是 datatypes.Date。
	// datatypes.Date 本质是 time.Time。
	// 让我们尝试显式转换类型，确保 GORM 知道我们在比较什么。
	// 或者直接使用字符串比较。
	//
	// 实际上，问题可能出在 TestGetSchedules 的 setup 并没有使用迁移后的数据库（它是内存数据库，且每次新建）。
	// 在 TestGetSchedules 中：
	// testDB.Model(&models.Schedule{}).Create(...)
	// 这里的 Create 使用的是 GORM + datatypes.Date。
	// datatypes.Date 在 Save 时会格式化为 "YYYY-MM-DD"。
	// 所以内存数据库里存的是 "2026-01-31"。
	//
	// 查询时：
	// Where("date >= ? AND date <= ?", "2026-01-01", "2026-01-31")
	// 数据库执行：SELECT * FROM schedules WHERE date >= '2026-01-01' AND date <= '2026-01-31'
	// 字符串比较 "2026-01-31" <= "2026-01-31" 是 TRUE。
	// 理论上应该能查到。
	//
	// 为什么测试失败？
	// schedule_handler_test.go:257: Expected 2 records, got 1. Dates: [2026-01-01]
	// 说明 01-31 没查到。
	//
	// 让我们调试一下，把 endDateStr 打印出来或者转换成 datatypes.Date 再传给 Where。

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

	// 开启事务
	tx := db.DB.Begin()

	for _, techID := range req.TechIDs {
		for _, dateStr := range req.Dates {
			// 解析日期 (验证格式)
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期格式错误: " + dateStr})
				return
			}
			// 转换为 datatypes.Date
			date := datatypes.Date(parsedDate)

			var schedule models.Schedule
			// 查找是否存在记录，存在则更新，不存在则创建 (Upsert)
			dbErr := tx.Where("tech_id = ? AND date = ?", techID, date).First(&schedule).Error

			if errors.Is(dbErr, gorm.ErrRecordNotFound) {
				// 创建新记录
				// Create using map to ensure zero values are respected
				if err := tx.Model(&models.Schedule{}).Create(map[string]interface{}{
					"tech_id":      techID,
					"date":         date,
					"is_available": req.IsAvailable,
					"time_slots":   datatypes.JSON([]byte("[]")),
				}).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建排班失败"})
					return
				}
			} else if dbErr != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询排班失败"})
				return
			} else {
				// 更新现有记录状态
				// 使用 Map 更新，确保 false 值也能被更新
				if err := tx.Model(&schedule).Updates(map[string]interface{}{
					"is_available": req.IsAvailable,
				}).Error; err != nil {
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
//   - unavailable: 不具备技能、休息中或忙碌的技师
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

	// 获取服务项目信息，计算结束时间
	var service models.ServiceProduct
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

	// ========== 新增：构建技师技能映射表 ==========
	// techID -> []serviceID
	techSkillsMap := make(map[uint][]uint)
	for _, tech := range allTechnicians {
		var techSkills []uint
		if tech.Skills != nil {
			// 尝试解析 JSON 数组 (新格式)
			if err := json.Unmarshal(tech.Skills, &techSkills); err != nil {
				// 如果解析失败，尝试解析旧格式（字符串数组）
				var oldSkills []string
				if err2 := json.Unmarshal(tech.Skills, &oldSkills); err2 == nil {
					// 尝试根据名称匹配服务项目ID
					for _, skillName := range oldSkills {
						var serviceItem models.ServiceProduct
						if err3 := db.DB.Where("name = ?", skillName).First(&serviceItem).Error; err3 == nil {
							techSkills = append(techSkills, serviceItem.ID)
						}
					}
				} else {
					techSkills = []uint{}
				}
			}
		}
		techSkillsMap[tech.ID] = techSkills
	}
	// ===========================================

	// 检查排班状态（是否在岗）
	checkDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	dateStr := checkDate.Format("2006-01-02")
	var schedules []models.Schedule
	db.DB.Where("date = ?", dateStr).Find(&schedules)

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
		techSkills := techSkillsMap[tech.ID]

		// ========== 新增：检查是否具备该服务项目的技能 ==========
		hasSkill := false
		for _, skillID := range techSkills {
			if skillID == service.ID {
				hasSkill = true
				break
			}
		}

		if !hasSkill {
			// 不具备该技能，加入不可用列表
			tech.Skills = nil // 避免JSON序列化问题
			tech.Reason = "skill_mismatch"
			unavailableTechnicians = append(unavailableTechnicians, tech)
			continue
		}
		// ======================================================

		// 检查是否在岗（默认在岗）
		isAvailable, exists := scheduleMap[tech.ID]
		if exists && !isAvailable {
			// 休息中
			tech.Skills = nil
			tech.Reason = "leave"
			unavailableTechnicians = append(unavailableTechnicians, tech)
			continue
		}

		// 检查是否有冲突预约
		if busyTechMap[tech.ID] {
			// 忙碌中
			tech.Skills = nil
			tech.Reason = "busy"
			unavailableTechnicians = append(unavailableTechnicians, tech)
			continue
		}

		// 可用
		tech.Skills = nil // 避免JSON序列化问题
		availableTechnicians = append(availableTechnicians, tech)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"available":   availableTechnicians,
			"unavailable": unavailableTechnicians,
			"total":       len(allTechnicians),
			"service": gin.H{ // 新增：返回服务信息供前端显示
				"id":   service.ID,
				"name": service.Name,
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
	var schedule models.Schedule
	// datatypes.Date可以直接用字符串查询
	scheduleErr := db.DB.Where("tech_id = ? AND date = ?", techIDStr, dateStr).First(&schedule).Error

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

	// 获取所有技师
	var allTechnicians []models.Technician
	if err := db.DB.Find(&allTechnicians).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Failed to fetch technicians"})
		return
	}

	// 1. 预处理：筛选具备技能的技师ID
	var skilledTechIDs []uint
	for _, tech := range allTechnicians {
		var techSkills []uint
		if tech.Skills != nil {
			// 尝试解析JSON
			json.Unmarshal(tech.Skills, &techSkills)
			// 如果是旧格式，这里简化处理，假设已经迁移或者只支持新格式。
			// 为了健壮性，这里简单复制之前的解析逻辑
			if len(techSkills) == 0 {
				var oldSkills []string
				if err2 := json.Unmarshal(tech.Skills, &oldSkills); err2 == nil {
					for _, skillName := range oldSkills {
						var serviceItem models.ServiceProduct
						if err3 := db.DB.Where("name = ?", skillName).First(&serviceItem).Error; err3 == nil {
							if serviceItem.ID == service.ID {
								techSkills = append(techSkills, serviceItem.ID)
								break
							}
						}
					}
				}
			}
		}

		for _, sID := range techSkills {
			if sID == service.ID {
				skilledTechIDs = append(skilledTechIDs, tech.ID)
				break
			}
		}
	}

	// 2. 预处理：获取当天排班（请假情况）
	// datatypes.Date 存储为 "YYYY-MM-DD"，查询时直接使用 dateStr 即可
	// checkDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var schedules []models.Schedule
	db.DB.Where("date = ?", dateStr).Find(&schedules)
	unavailableScheduleTechs := make(map[uint]bool)
	for _, s := range schedules {
		if !s.IsAvailable {
			unavailableScheduleTechs[s.TechID] = true
		}
	}

	// 3. 预处理：获取当天所有预约
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	var appointments []models.Appointment
	db.DB.Where("start_time >= ? AND start_time < ? AND status = ?", startOfDay, endOfDay, "pending").
		Find(&appointments)

	// 生成时间槽 (10:00 - 22:00, 30分钟间隔)
	// TODO: 营业时间应从配置读取
	openTime := time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, time.UTC)
	closeTime := time.Date(date.Year(), date.Month(), date.Day(), 22, 0, 0, 0, time.UTC)

	type TimeSlot struct {
		Time           string `json:"time"`            // "10:00"
		Status         string `json:"status"`          // "available", "waitlist", "closed"
		AvailableCount int    `json:"available_count"` // 可用技师数
		StartTime      string `json:"start_time"`      // ISO String
	}

	var slots []TimeSlot

	for t := openTime; t.Before(closeTime); t = t.Add(30 * time.Minute) {
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
