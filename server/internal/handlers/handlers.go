package handlers

import (
	"fmt"
	"net/http"
	"time"

	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"
	"smartspa-admin/internal/response"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats 获取仪表盘统计数据
func GetDashboardStats(c *gin.Context) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	// 1. 今日营收（已完成订单的实付金额总和）
	var dailyRevenue float64
	if err := db.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = DATE(?)", "completed", today).
		Select("COALESCE(SUM(actual_paid), 0)").
		Scan(&dailyRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to calculate daily revenue", err))
		return
	}

	// 昨日营收（用于计算增长率）
	var yesterdayRevenue float64
	if err := db.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = DATE(?)", "completed", yesterday).
		Select("COALESCE(SUM(actual_paid), 0)").
		Scan(&yesterdayRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to calculate yesterday revenue", err))
		return
	}

	// 计算增长率
	var revenueGrowth float64
	if yesterdayRevenue > 0 {
		revenueGrowth = ((dailyRevenue - yesterdayRevenue) / yesterdayRevenue) * 100
	}

	// 2. 今日新增会员
	var newMembers int64
	if err := db.DB.Model(&models.Member{}).
		Where("DATE(created_at) = DATE(?)", today).
		Count(&newMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count new members", err))
		return
	}

	// 3. 活跃技师数量（状态为 free 或 booked 的技师）
	var activeTechs int64
	if err := db.DB.Model(&models.Technician{}).
		Where("status IN ?", []int{0, 1}).
		Count(&activeTechs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count active techs", err))
		return
	}

	// 总技师数
	var totalTechs int64
	if err := db.DB.Model(&models.Technician{}).
		Count(&totalTechs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count total techs", err))
		return
	}

	// 4. 技师负载率（今日已完成预约数 / 总技师数 / 8小时 * 100）
	var todayCompletedAppointments int64
	if err := db.DB.Model(&models.Appointment{}).
		Where("status = ? AND DATE(start_time) = DATE(?)", "completed", today).
		Count(&todayCompletedAppointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count today appointments", err))
		return
	}

	var occupancyRate float64
	if totalTechs > 0 {
		// 假设每个技师每天工作8小时，平均每个服务1.5小时
		maxCapacity := float64(totalTechs) * 8 / 1.5
		occupancyRate = (float64(todayCompletedAppointments) / maxCapacity) * 100
		if occupancyRate > 100 {
			occupancyRate = 100
		}
	}

	stats := gin.H{
		"dailyRevenue":  dailyRevenue,
		"revenueGrowth": revenueGrowth,
		"newMembers":    newMembers,
		"activeTechs":   activeTechs,
		"occupancyRate": occupancyRate,
	}

	c.JSON(http.StatusOK, response.Success(stats, ""))
}

// ListAppointments 获取订单列表 (支持状态筛选)
func ListAppointments(c *gin.Context) {
	var appointments []models.Appointment
	query := db.DB.Model(&models.Appointment{})

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.
		Preload("Member").
		Preload("Technician").
		Preload("ServiceItem").
		Limit(50).
		Order("created_at DESC").
		Find(&appointments).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(appointments, ""))
}

// GetFissionRanking 获取分销排行榜
func GetFissionRanking(c *gin.Context) {
	type FissionRank struct {
		ID              uint    `json:"id"`
		Name            string  `json:"name"`
		Phone           string  `json:"phone"`
		Level           string  `json:"level"`
		InviteCount     int64   `json:"inviteCount"`
		TotalCommission float64 `json:"totalCommission"`
	}

	var rankings []FissionRank

	// 统计每个会员邀请的人数和累计佣金
	if err := db.DB.Model(&models.Member{}).
		Select("members.id, members.name, members.phone, members.level, COUNT(invitees.id) as invite_count, COALESCE(SUM(fission_logs.commission_amount), 0) as total_commission").
		Joins("LEFT JOIN members AS invitees ON invitees.referrer_id = members.id").
		Joins("LEFT JOIN fission_logs ON fission_logs.inviter_id = members.id").
		Group("members.id, members.name, members.phone, members.level").
		Having("COUNT(invitees.id) > 0").
		Order("invite_count DESC, total_commission DESC").
		Limit(10).
		Scan(&rankings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to get fission ranking", err))
		return
	}

	c.JSON(http.StatusOK, response.Success(rankings, ""))
}

// CreateAppointment 创建预约
func CreateAppointment(c *gin.Context) {
	var req struct {
		MemberID      uint   `json:"member_id"`
		TechID        uint   `json:"tech_id"`
		ServiceID     uint   `json:"service_id"`
		StartTime     string `json:"start_time"`
		AllowWaitlist bool   `json:"allow_waitlist"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Parse StartTime
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid time format", nil))
		return
	}

	// Get Service details for duration and price
	var service models.ServiceProduct
	if err := db.DB.First(&service, req.ServiceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Service not found", nil))
		return
	}

	endTime := startTime.Add(time.Duration(service.Duration) * time.Minute)

	// Check Schedule Availability
	// Normalize date to UTC midnight for comparison with Schedule table
	checkDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	var schedule models.Schedule
	if err := db.DB.Where("tech_id = ? AND date = ?", req.TechID, checkDate).First(&schedule).Error; err == nil {
		if !schedule.IsAvailable {
			c.JSON(http.StatusConflict, response.Error(http.StatusConflict, "Technician is on leave/unavailable on this date", nil))
			return
		}
	}

	// Conflict Detection
	var conflictCount int64
	db.DB.Model(&models.Appointment{}).
		Where("tech_id = ? AND status = 'pending' AND start_time < ? AND end_time > ?",
			req.TechID, endTime, startTime).
		Count(&conflictCount)

	status := "pending"
	if conflictCount > 0 {
		if !req.AllowWaitlist {
			c.JSON(http.StatusConflict, response.Error(http.StatusConflict, "Technician is busy at this time", nil))
			return
		}
		status = "waiting"
	}

	appointment := models.Appointment{
		MemberID:    req.MemberID,
		TechID:      req.TechID,
		ServiceID:   req.ServiceID,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      status,
		OriginPrice: service.Price,
		ActualPrice: service.Price,
	}

	if err := db.DB.Create(&appointment).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Reload to get associations
	db.DB.Preload("Member").Preload("Technician").Preload("ServiceItem").First(&appointment, appointment.ID)

	msg := ""
	if status == "waiting" {
		msg = "Time slot conflict. Added to waitlist."
	}

	c.JSON(http.StatusOK, response.Success(appointment, msg))
}

// CancelAppointment 取消预约并触发候补检查
func CancelAppointment(c *gin.Context) {
	id := c.Param("id")
	var appt models.Appointment
	if err := db.DB.First(&appt, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Appointment not found", nil))
		return
	}

	if appt.Status == "cancelled" {
		c.JSON(http.StatusOK, response.Success(nil, "Already cancelled"))
		return
	}

	appt.Status = "cancelled"
	if err := db.DB.Save(&appt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to cancel appointment", nil))
		return
	}

	// Trigger Waitlist Check for this technician
	go checkWaitlist(appt.TechID)

	c.JSON(http.StatusOK, response.Success(nil, "Appointment cancelled"))
}

// checkWaitlist checks if any waiting appointments can be promoted
func checkWaitlist(techID uint) {
	var waitingList []models.Appointment
	// Find all waiting appointments for this tech, ordered by creation time (FCFS)
	db.DB.Where("tech_id = ? AND status = ?", techID, "waiting").
		Order("created_at asc").
		Find(&waitingList)

	for _, waitAppt := range waitingList {
		// Check if this slot is now free
		var conflictCount int64
		db.DB.Model(&models.Appointment{}).
			Where("tech_id = ? AND status = 'pending' AND start_time < ? AND end_time > ?",
				techID, waitAppt.EndTime, waitAppt.StartTime).
			Count(&conflictCount)

		if conflictCount == 0 {
			// Promote to pending
			waitAppt.Status = "pending"
			db.DB.Save(&waitAppt)
		}
	}
}

// CompleteAppointment 完成预约并结算
func CompleteAppointment(c *gin.Context) {
	id := c.Param("id")
	var appt models.Appointment
	if err := db.DB.Preload("Member").First(&appt, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Appointment not found", nil))
		return
	}

	if appt.Status == "completed" {
		c.JSON(http.StatusOK, response.Success(nil, "Already completed"))
		return
	}

	// Start Transaction
	tx := db.DB.Begin()

	// 1. Update Appointment Status
	appt.Status = "completed"
	if err := tx.Save(&appt).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update appointment", nil))
		return
	}

	// 2. Update Member Consumption & Level
	member := appt.Member
	member.YearlyTotalConsumption += appt.ActualPrice

	// Simple Level Logic
	if member.YearlyTotalConsumption > 10000 {
		member.Level = "platinum"
	} else if member.YearlyTotalConsumption > 5000 {
		member.Level = "gold"
	} else if member.YearlyTotalConsumption > 1000 {
		member.Level = "silver"
	}

	if err := tx.Save(&member).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update member", nil))
		return
	}

	// 3. Commission Logic
	if member.ReferrerID != nil {
		commission := appt.ActualPrice * 0.10 // 10%

		// Update Referrer Balance
		var referrer models.Member
		if err := tx.First(&referrer, *member.ReferrerID).Error; err == nil {
			referrer.Balance += commission
			if err := tx.Save(&referrer).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update referrer", nil))
				return
			}

			// Log Fission
			fissionLog := models.FissionLog{
				InviterID:        referrer.ID,
				InviteeID:        member.ID,
				CommissionAmount: commission,
				OrderID:          &appt.ID,
			}
			if err := tx.Create(&fissionLog).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create fission log", nil))
				return
			}
		}
	}

	tx.Commit()

	// Trigger Waitlist Check for this technician
	go checkWaitlist(appt.TechID)

	c.JSON(http.StatusOK, response.Success(nil, "Appointment completed and settled"))
}

// ListTechnicians 获取技师列表
func ListTechnicians(c *gin.Context) {
	var technicians []models.Technician
	if err := db.DB.Find(&technicians).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(technicians, ""))
}

// CreateTechnician 创建技师
func CreateTechnician(c *gin.Context) {
	var tech models.Technician
	if err := c.ShouldBindJSON(&tech); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if tech.Name == "" {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Name is required", nil))
		return
	}

	if err := db.DB.Create(&tech).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create technician", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(tech, "Technician created successfully"))
}

// UpdateTechnician 更新技师信息
func UpdateTechnician(c *gin.Context) {
	id := c.Param("id")
	var tech models.Technician
	if err := db.DB.First(&tech, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Technician not found", nil))
		return
	}

	var req models.Technician
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	tech.Name = req.Name
	tech.Skills = req.Skills
	tech.Status = req.Status

	if err := db.DB.Save(&tech).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update technician", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(tech, "Technician updated successfully"))
}

// DeleteTechnician 删除技师
func DeleteTechnician(c *gin.Context) {
	id := c.Param("id")

	// 开启事务
	tx := db.DB.Begin()

	// 查找该技师是否有待服务的订单（status = 'pending'）
	var pendingAppointments []models.Appointment
	if err := tx.Where("tech_id = ? AND status = ?", id, "pending").Find(&pendingAppointments).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to query appointments", nil))
		return
	}

	// 如果有待服务的订单，将状态修改为候补中（waitlist）
	if len(pendingAppointments) > 0 {
		if err := tx.Model(&models.Appointment{}).Where("tech_id = ? AND status = ?", id, "pending").Update("status", "waitlist").Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update appointments to waitlist", nil))
			return
		}
	}

	// 删除技师
	if err := tx.Delete(&models.Technician{}, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to delete technician", nil))
		return
	}

	tx.Commit()

	msg := "Technician deleted successfully"
	if len(pendingAppointments) > 0 {
		msg = fmt.Sprintf("%s. %d pending appointments moved to waitlist", msg, len(pendingAppointments))
	}

	c.JSON(http.StatusOK, response.Success(nil, msg))
}

// ListServiceItems 获取服务项目
func ListServiceItems(c *gin.Context) {
	var items []models.ServiceProduct
	query := db.DB.Model(&models.ServiceProduct{})

	if c.Query("active_only") == "true" {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(items, ""))
}

// CreateServiceItem 创建服务项目
func CreateServiceItem(c *gin.Context) {
	var item models.ServiceProduct
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if item.Name == "" || item.Price <= 0 || item.Duration <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid service item data", nil))
		return
	}

	if err := db.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create service item", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(item, "Service item created successfully"))
}

// UpdateServiceItem 更新服务项目
func UpdateServiceItem(c *gin.Context) {
	id := c.Param("id")
	var item models.ServiceProduct
	if err := db.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Service item not found", nil))
		return
	}

	var req models.ServiceProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	item.Name = req.Name
	item.Duration = req.Duration
	item.Price = req.Price
	item.IsActive = req.IsActive
	item.ImageURL = req.ImageURL

	if err := db.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update service item", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(item, "Service item updated successfully"))
}

// DeleteServiceItem 删除服务项目
func DeleteServiceItem(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.ServiceProduct{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to delete service item", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil, "Service item deleted successfully"))
}

// ListMembers 获取会员列表
func ListMembers(c *gin.Context) {
	var members []models.Member
	if err := db.DB.Limit(20).Find(&members).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(members, ""))
}

// CreateMember 创建会员
func CreateMember(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if member.Name == "" || member.Phone == "" {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Name and Phone are required", nil))
		return
	}

	if err := db.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create member", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(member, "Member created successfully"))
}
