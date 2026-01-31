package handlers

import (
	"net/http"
	"strconv"

	"server/internal/db"
	"server/internal/models"
	"server/internal/response"
	"server/pkg/config"
	"server/pkg/util"

	"github.com/gin-gonic/gin"
)

// InventoryChangeRequest represents the request body for inventory changes
type InventoryChangeRequest struct {
	ProductID    uint     `json:"product_id" binding:"required"`
	ChangeAmount int      `json:"change_amount" binding:"required"`
	ActionType   string   `json:"action_type" binding:"required,oneof=restock sale adjustment"`
	MemberID     *uint    `json:"member_id"`   // 购买者ID（销售时可选）
	SaleAmount   *float64 `json:"sale_amount"` // 销售金额（销售时可选）
	Remark       string   `json:"remark"`
}

// ListInventoryLogs returns all inventory logs
func ListInventoryLogs(c *gin.Context) {
	database := db.GetDB()

	query := database.Preload("Product").Preload("Operator").Order("created_at DESC")

	// Filter by product ID if specified
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		productID, err := strconv.ParseUint(productIDStr, 10, 32)
		if err == nil {
			query = query.Where("product_id = ?", productID)
		}
	}

	// Filter by action type if specified
	if actionType := c.Query("action_type"); actionType != "" {
		query = query.Where("action_type = ?", actionType)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.InventoryLog{}).Count(&total)

	var logs []models.InventoryLog
	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch inventory logs", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, ""))
}

// GetProductInventoryLogs returns inventory logs for a specific product
func GetProductInventoryLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid product ID", nil))
		return
	}

	database := db.GetDB()

	var logs []models.InventoryLog
	if err := database.Where("product_id = ?", id).
		Preload("Operator").
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch inventory logs", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"logs":  logs,
		"total": len(logs),
	}, ""))
}

// CreateInventoryChange creates a new inventory change record
func CreateInventoryChange(c *gin.Context) {
	var req InventoryChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid request: "+err.Error(), nil))
		return
	}

	// Get current user ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "User not authenticated", nil))
		return
	}

	database := db.GetDB()

	// Validate action type and change amount
	if req.ActionType == "restock" && req.ChangeAmount <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Restock amount must be positive", nil))
		return
	}
	if req.ActionType == "sale" && req.ChangeAmount >= 0 {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Sale amount must be negative", nil))
		return
	}

	// Use transaction to ensure atomicity
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get product and lock for update
	var product models.PhysicalProduct
	if err := tx.Clauses().First(&product, req.ProductID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Product not found", nil))
		return
	}

	// Check if stock is sufficient for sale or negative adjustment
	if req.ChangeAmount < 0 && product.Stock+req.ChangeAmount < 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Insufficient stock", gin.H{
			"current_stock":    product.Stock,
			"requested_change": req.ChangeAmount,
		}))
		return
	}

	beforeStock := product.Stock
	afterStock := beforeStock + req.ChangeAmount

	// Update product stock
	product.Stock = afterStock
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update product stock", nil))
		return
	}

	// Create inventory log
	inventoryLog := models.InventoryLog{
		ProductID:    req.ProductID,
		OperatorID:   userID.(uint),
		MemberID:     req.MemberID,
		ChangeAmount: req.ChangeAmount,
		ActionType:   req.ActionType,
		BeforeStock:  beforeStock,
		AfterStock:   afterStock,
		SaleAmount:   req.SaleAmount,
		Remark:       req.Remark,
	}

	if err := tx.Create(&inventoryLog).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create inventory log", nil))
		return
	}

	// Fission commission logic for product sales with member
	if req.ActionType == "sale" && req.MemberID != nil && req.SaleAmount != nil && *req.SaleAmount > 0 {
		var member models.Member
		if err := tx.First(&member, *req.MemberID).Error; err == nil && member.ReferrerID != nil {
			commissionAmount := util.CalculateRate(*req.SaleAmount, config.GlobalCommission.ReferralRate)
			commissionInCents := util.ToCents(commissionAmount)
			saleAmountInCents := util.ToCents(*req.SaleAmount)

			if commissionInCents >= 0 && commissionInCents <= saleAmountInCents {
				var referrer models.Member
				if err := tx.First(&referrer, *member.ReferrerID).Error; err == nil {
					oldBalance := referrer.Balance
					referrer.Balance += float64(commissionInCents) / 100

					if referrer.Balance >= 0 {
						expectedChange := float64(commissionInCents) / 100
						actualChange := referrer.Balance - oldBalance
						if actualChange >= expectedChange-0.01 && actualChange <= expectedChange+0.01 {
							if err := tx.Save(&referrer).Error; err == nil {
								fissionLog := models.FissionLog{
									InviterID:        referrer.ID,
									InviteeID:        member.ID,
									CommissionAmount: float64(commissionInCents) / 100,
								}
								tx.Create(&fissionLog)
							}
						}
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to commit transaction", nil))
		return
	}

	// Reload log with associations
	database.Preload("Product").Preload("Operator").First(&inventoryLog, inventoryLog.ID)

	c.JSON(http.StatusOK, response.Success(inventoryLog, "Inventory updated successfully"))
}

// BatchRestockRequest represents the request body for batch restocking
type BatchRestockRequest struct {
	Items []struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,min=1"`
		Remark    string `json:"remark"`
	} `json:"items" binding:"required,min=1"`
}

// BatchRestock creates multiple restock records at once
func BatchRestock(c *gin.Context) {
	var req BatchRestockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid request: "+err.Error(), nil))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "User not authenticated", nil))
		return
	}

	database := db.GetDB()
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var createdLogs []models.InventoryLog

	for _, item := range req.Items {
		// Get product
		var product models.PhysicalProduct
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Product not found: ID "+strconv.Itoa(int(item.ProductID)), nil))
			return
		}

		beforeStock := product.Stock
		afterStock := beforeStock + item.Quantity

		// Update stock
		product.Stock = afterStock
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update product stock", nil))
			return
		}

		// Create log
		log := models.InventoryLog{
			ProductID:    item.ProductID,
			OperatorID:   userID.(uint),
			ChangeAmount: item.Quantity,
			ActionType:   "restock",
			BeforeStock:  beforeStock,
			AfterStock:   afterStock,
			Remark:       item.Remark,
		}
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create inventory log", nil))
			return
		}

		createdLogs = append(createdLogs, log)
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to commit transaction", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"success_count": len(createdLogs),
		"logs":          createdLogs,
	}, "Batch restock completed successfully"))
}

// GetInventoryStats returns inventory statistics
func GetInventoryStats(c *gin.Context) {
	database := db.GetDB()

	var stats struct {
		TotalTransactions int64 `json:"total_transactions"`
		RestockCount      int64 `json:"restock_count"`
		SaleCount         int64 `json:"sale_count"`
		AdjustmentCount   int64 `json:"adjustment_count"`
	}

	// Total transactions
	database.Model(&models.InventoryLog{}).Count(&stats.TotalTransactions)

	// Count by action type
	database.Model(&models.InventoryLog{}).Where("action_type = ?", "restock").Count(&stats.RestockCount)
	database.Model(&models.InventoryLog{}).Where("action_type = ?", "sale").Count(&stats.SaleCount)
	database.Model(&models.InventoryLog{}).Where("action_type = ?", "adjustment").Count(&stats.AdjustmentCount)

	c.JSON(http.StatusOK, response.Success(stats, ""))
}
