package handlers

import (
	"net/http"
	"strconv"

	"server/internal/db"
	"server/internal/models"
	"server/internal/response"

	"github.com/gin-gonic/gin"
)

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	RetailPrice float64 `json:"retail_price" binding:"required,min=0"`
	CostPrice   float64 `json:"cost_price" binding:"required,min=0"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	RetailPrice float64 `json:"retail_price" binding:"min=0"`
	CostPrice   float64 `json:"cost_price" binding:"min=0"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

// ListProducts returns all physical products
func ListProducts(c *gin.Context) {
	var products []models.PhysicalProduct
	database := db.GetDB()

	query := database.Order("created_at DESC")

	// Filter by active status if specified
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch products", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"products": products,
		"total":    len(products),
	}, ""))
}

// GetProduct returns a single product by ID
func GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid product ID", nil))
		return
	}

	var product models.PhysicalProduct
	database := db.GetDB()

	if err := database.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Product not found", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(product, ""))
}

// CreateProduct creates a new physical product
func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid request: "+err.Error(), nil))
		return
	}

	// Get current user ID for the initial inventory log
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "User not authenticated", nil))
		return
	}

	database := db.GetDB()

	// Create product
	product := models.PhysicalProduct{
		Name:        req.Name,
		Stock:       req.Stock,
		RetailPrice: req.RetailPrice,
		CostPrice:   req.CostPrice,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
	}

	// Use transaction to ensure product and initial inventory log are created together
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create product", nil))
		return
	}

	// Create initial inventory log if stock > 0
	if req.Stock > 0 {
		inventoryLog := models.InventoryLog{
			ProductID:    product.ID,
			OperatorID:   userID.(uint),
			ChangeAmount: req.Stock,
			ActionType:   "restock",
			BeforeStock:  0,
			AfterStock:   req.Stock,
			Remark:       "Initial stock",
		}
		if err := tx.Create(&inventoryLog).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create inventory log", nil))
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to commit transaction", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(product, "Product created successfully"))
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid product ID", nil))
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid request: "+err.Error(), nil))
		return
	}

	database := db.GetDB()

	var product models.PhysicalProduct
	if err := database.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Product not found", nil))
		return
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.RetailPrice > 0 {
		product.RetailPrice = req.RetailPrice
	}
	if req.CostPrice > 0 {
		product.CostPrice = req.CostPrice
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := database.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to update product", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(product, "Product updated successfully"))
}

// DeleteProduct soft deletes a product
func DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid product ID", nil))
		return
	}

	database := db.GetDB()

	var product models.PhysicalProduct
	if err := database.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Product not found", nil))
		return
	}

	// Soft delete
	if err := database.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to delete product", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil, "Product deleted successfully"))
}

// GetProductStats returns statistics about products
func GetProductStats(c *gin.Context) {
	database := db.GetDB()

	var stats struct {
		TotalProducts   int64   `json:"total_products"`
		ActiveProducts  int64   `json:"active_products"`
		TotalValue      float64 `json:"total_value"`        // 库存总价值（按零售价）
		LowStockCount   int64   `json:"low_stock_count"`    // 低库存商品数
		OutOfStockCount int64   `json:"out_of_stock_count"` // 零库存商品数
	}

	// Total products
	database.Model(&models.PhysicalProduct{}).Count(&stats.TotalProducts)

	// Active products
	database.Model(&models.PhysicalProduct{}).Where("is_active = ?", true).Count(&stats.ActiveProducts)

	// Total inventory value
	var products []models.PhysicalProduct
	database.Find(&products)
	for _, p := range products {
		stats.TotalValue += float64(p.Stock) * p.RetailPrice
	}

	// Low stock (< 10) and out of stock
	database.Model(&models.PhysicalProduct{}).Where("stock < ? AND stock > 0", 10).Count(&stats.LowStockCount)
	database.Model(&models.PhysicalProduct{}).Where("stock = 0").Count(&stats.OutOfStockCount)

	c.JSON(http.StatusOK, response.Success(stats, ""))
}
