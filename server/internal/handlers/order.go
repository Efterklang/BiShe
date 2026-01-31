package handlers

import (
	"net/http"
	"strconv"
	"time"

	"server/internal/db"
	"server/internal/models"
	"server/internal/response"
	"server/pkg/config"
	"server/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	OrderType      string `json:"order_type" binding:"required,oneof=service physical"`
	AppointmentID  *uint  `json:"appointment_id"`
	InventoryLogID *uint  `json:"inventory_log_id"`
}

func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid request: "+err.Error(), nil))
		return
	}

	if req.OrderType == "service" && (req.AppointmentID == nil || req.InventoryLogID != nil) {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Service order requires appointment_id and forbids inventory_log_id", nil))
		return
	}
	if req.OrderType == "physical" && (req.InventoryLogID == nil || req.AppointmentID != nil) {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Physical order requires inventory_log_id and forbids appointment_id", nil))
		return
	}

	database := db.GetDB()
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existing models.Order
	if req.OrderType == "service" {
		if err := tx.Where("appointment_id = ?", *req.AppointmentID).First(&existing).Error; err == nil {
			tx.Commit()
			c.JSON(http.StatusOK, response.Success(existing, ""))
			return
		} else if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to check existing order", err.Error()))
			return
		}

		var appt models.Appointment
		if err := tx.Preload("Member").First(&appt, *req.AppointmentID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Appointment not found", err.Error()))
			return
		}
		if appt.Status != "completed" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Appointment not completed", nil))
			return
		}

		paidAmount := appt.ActualPrice
		inviterID := appt.Member.ReferrerID
		commissionAmount := 0.0
		if inviterID != nil {
			commissionAmount = util.CalculateRate(paidAmount, config.GlobalCommission.ReferralRate)
		}

		order := models.Order{
			MemberID:         appt.MemberID,
			InviterID:        inviterID,
			PaidAmount:       paidAmount,
			CommissionAmount: commissionAmount,
			OrderType:        "service",
			AppointmentID:    req.AppointmentID,
		}
		order.CreatedAt = appt.EndTime
		order.UpdatedAt = appt.EndTime
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Failed to create order", err.Error()))
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, response.Success(order, ""))
		return
	}

	if err := tx.Where("inventory_log_id = ?", *req.InventoryLogID).First(&existing).Error; err == nil {
		tx.Commit()
		c.JSON(http.StatusOK, response.Success(existing, ""))
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to check existing order", err.Error()))
		return
	}

	var invLog models.InventoryLog
	if err := tx.Preload("Product").Preload("Member").First(&invLog, *req.InventoryLogID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Inventory log not found", err.Error()))
		return
	}
	if invLog.ActionType != "sale" || invLog.ChangeAmount >= 0 || invLog.MemberID == nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Inventory log is not a valid sale", nil))
		return
	}

	paidAmount := 0.0
	if invLog.SaleAmount != nil {
		paidAmount = *invLog.SaleAmount
	} else {
		paidAmount = float64(-invLog.ChangeAmount) * invLog.Product.RetailPrice
	}

	var inviterID *uint
	if invLog.Member != nil {
		inviterID = invLog.Member.ReferrerID
	}

	commissionAmount := 0.0
	if inviterID != nil {
		commissionAmount = util.CalculateRate(paidAmount, config.GlobalCommission.ReferralRate)
	}

	order := models.Order{
		MemberID:         *invLog.MemberID,
		InviterID:        inviterID,
		PaidAmount:       paidAmount,
		CommissionAmount: commissionAmount,
		OrderType:        "physical",
		InventoryLogID:   req.InventoryLogID,
	}
	order.CreatedAt = invLog.CreatedAt
	order.UpdatedAt = invLog.CreatedAt
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Failed to create order", err.Error()))
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, response.Success(order, ""))
}

func ListOrders(c *gin.Context) {
	database := db.GetDB()

	query := database.Model(&models.Order{}).Preload("Member").Preload("Inviter").Order("created_at DESC")

	if orderType := c.Query("order_type"); orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}
	if memberIDStr := c.Query("member_id"); memberIDStr != "" {
		if memberID, err := strconv.ParseUint(memberIDStr, 10, 32); err == nil {
			query = query.Where("member_id = ?", uint(memberID))
		}
	}
	if memberLevel := c.Query("member_level"); memberLevel != "" {
		query = query.Joins("JOIN members ON members.id = orders.member_id").Where("members.level = ?", memberLevel)
	}

	start, end := parseTimeRange(c.Query("start"), c.Query("end"))
	if !start.IsZero() {
		query = query.Where("orders.created_at >= ?", start)
	}
	if !end.IsZero() {
		query = query.Where("orders.created_at < ?", end)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count orders", err.Error()))
		return
	}

	var orders []models.Order
	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch orders", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"orders":    orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, ""))
}

func parseTimeRange(startStr, endStr string) (time.Time, time.Time) {
	parse := func(s string) time.Time {
		if s == "" {
			return time.Time{}
		}
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t
		}
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			return t
		}
		return time.Time{}
	}

	start := parse(startStr)
	end := parse(endStr)
	if !end.IsZero() && end.Hour() == 0 && end.Minute() == 0 && end.Second() == 0 && end.Nanosecond() == 0 {
		end = end.Add(24 * time.Hour)
	}
	return start, end
}
