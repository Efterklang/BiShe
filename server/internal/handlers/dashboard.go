package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"smartspa-admin/internal/models"
	"smartspa-admin/internal/response"
)

// DashboardHandler handles dashboard-related requests
type DashboardHandler struct {
	db *gorm.DB
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetStats returns dashboard statistics
// GET /api/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	// 1. 今日营收（已完成订单的实付金额总和）
	var dailyRevenue float64
	if err := h.db.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = DATE(?)", "completed", today).
		Select("COALESCE(SUM(actual_paid), 0)").
		Scan(&dailyRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to calculate daily revenue", err.Error()))
		return
	}

	// 昨日营收（用于计算增长率）
	var yesterdayRevenue float64
	if err := h.db.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = DATE(?)", "completed", yesterday).
		Select("COALESCE(SUM(actual_paid), 0)").
		Scan(&yesterdayRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to calculate yesterday revenue", err.Error()))
		return
	}

	// 计算增长率
	var revenueGrowth float64
	if yesterdayRevenue > 0 {
		revenueGrowth = ((dailyRevenue - yesterdayRevenue) / yesterdayRevenue) * 100
	}

	// 2. 今日新增会员
	var newMembers int64
	if err := h.db.Model(&models.Member{}).
		Where("DATE(created_at) = DATE(?)", today).
		Count(&newMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count new members", err.Error()))
		return
	}

	// 3. 活跃技师数量（状态为 free 或 booked 的技师）
	var activeTechs int64
	if err := h.db.Model(&models.Technician{}).
		Where("status IN ?", []int{0, 1}).
		Count(&activeTechs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count active techs", err.Error()))
		return
	}

	// 总技师数
	var totalTechs int64
	if err := h.db.Model(&models.Technician{}).
		Count(&totalTechs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count total techs", err.Error()))
		return
	}

	// 4. 技师负载率（今日已完成预约数 / 总技师数 / 8小时 * 100）
	var todayCompletedAppointments int64
	if err := h.db.Model(&models.Appointment{}).
		Where("status = ? AND DATE(start_time) = DATE(?)", "completed", today).
		Count(&todayCompletedAppointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count today appointments", err.Error()))
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

	c.JSON(http.StatusOK, response.Success(gin.H{
		"dailyRevenue":  dailyRevenue,
		"revenueGrowth": revenueGrowth,
		"newMembers":    newMembers,
		"activeTechs":   activeTechs,
		"occupancyRate": occupancyRate,
		"totalTechs":    totalTechs,
		"todayOrders":   todayCompletedAppointments,
	}, ""))
}

// GetRevenueTrend returns revenue trend for a specified period
// GET /api/dashboard/revenue-trend?days=7|30|90
func (h *DashboardHandler) GetRevenueTrend(c *gin.Context) {
	// 获取时间范围参数，默认30天
	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
			if days > 365 {
				days = 365 // 限制最多一年
			}
		}
	}

	now := time.Now()
	startDate := now.AddDate(0, 0, -days)

	type DailyRevenue struct {
		Date           string  `json:"date"`
		ServiceRevenue float64 `json:"service_revenue"`
		ProductRevenue float64 `json:"product_revenue"`
	}

	// 统计服务营收（从 appointments 表）
	var serviceRevenues []struct {
		Date    string  `json:"date"`
		Revenue float64 `json:"revenue"`
	}
	if err := h.db.Model(&models.Appointment{}).
		Select("DATE(created_at) as date, COALESCE(SUM(actual_price), 0) as revenue").
		Where("status = ? AND created_at >= ?", "completed", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&serviceRevenues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to get service revenue trend", err.Error()))
		return
	}

	// 统计商品营收（从 orders 表，type='product'）
	var productRevenues []struct {
		Date    string  `json:"date"`
		Revenue float64 `json:"revenue"`
	}
	if err := h.db.Model(&models.Order{}).
		Select("DATE(created_at) as date, COALESCE(SUM(actual_paid), 0) as revenue").
		Where("type = ? AND status = ? AND created_at >= ?", "product", "completed", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&productRevenues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to get product revenue trend", err.Error()))
		return
	}

	// 填充缺失的日期
	serviceDateMap := make(map[string]float64)
	for _, sr := range serviceRevenues {
		serviceDateMap[sr.Date] = sr.Revenue
	}

	productDateMap := make(map[string]float64)
	for _, pr := range productRevenues {
		productDateMap[pr.Date] = pr.Revenue
	}

	result := make([]DailyRevenue, days)
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -(days-1)+i).Format("2006-01-02")
		result[i] = DailyRevenue{
			Date:           date,
			ServiceRevenue: serviceDateMap[date],
			ProductRevenue: productDateMap[date],
		}
	}

	c.JSON(http.StatusOK, response.Success(result, ""))
}

// GetServiceRanking returns top popular services
// GET /api/dashboard/service-ranking
func (h *DashboardHandler) GetServiceRanking(c *gin.Context) {
	type ServiceRank struct {
		ServiceID    uint    `json:"service_id"`
		ServiceName  string  `json:"service_name"`
		OrderCount   int64   `json:"order_count"`
		TotalRevenue float64 `json:"total_revenue"`
	}

	var rankings []ServiceRank

	// 统计近30天各服务的订单数和营收
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	if err := h.db.Model(&models.Appointment{}).
		Select("service_products.id as service_id, service_products.name as service_name, COUNT(appointments.id) as order_count, COALESCE(SUM(appointments.actual_price), 0) as total_revenue").
		Joins("JOIN service_products ON service_products.id = appointments.service_id").
		Where("appointments.status = ? AND appointments.created_at >= ?", "completed", thirtyDaysAgo).
		Group("service_products.id, service_products.name").
		Order("order_count DESC").
		Limit(10).
		Scan(&rankings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to get service ranking", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(rankings, ""))
}

// GetFissionRanking returns top members by invitation count
// GET /api/fission/ranking
func (h *DashboardHandler) GetFissionRanking(c *gin.Context) {
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
	if err := h.db.Model(&models.Member{}).
		Select("members.id, members.name, members.phone, members.level, COUNT(invitees.id) as invite_count, COALESCE(SUM(fission_logs.commission_amount), 0) as total_commission").
		Joins("LEFT JOIN members AS invitees ON invitees.referrer_id = members.id").
		Joins("LEFT JOIN fission_logs ON fission_logs.inviter_id = members.id").
		Group("members.id, members.name, members.phone, members.level").
		Having("COUNT(invitees.id) > 0").
		Order("invite_count DESC, total_commission DESC").
		Limit(10).
		Scan(&rankings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to get fission ranking", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(rankings, ""))
}

// GetMonthlyStats returns monthly statistics
// GET /api/dashboard/monthly-stats
func (h *DashboardHandler) GetMonthlyStats(c *gin.Context) {
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 本月营收
	var monthlyRevenue float64
	if err := h.db.Model(&models.Order{}).
		Where("status = ? AND created_at >= ?", "completed", firstDayOfMonth).
		Select("COALESCE(SUM(actual_paid), 0)").
		Scan(&monthlyRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to calculate monthly revenue", err.Error()))
		return
	}

	// 本月新增会员
	var monthlyNewMembers int64
	if err := h.db.Model(&models.Member{}).
		Where("created_at >= ?", firstDayOfMonth).
		Count(&monthlyNewMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count monthly new members", err.Error()))
		return
	}

	// 本月完成订单数
	var monthlyOrders int64
	if err := h.db.Model(&models.Order{}).
		Where("status = ? AND created_at >= ?", "completed", firstDayOfMonth).
		Count(&monthlyOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count monthly orders", err.Error()))
		return
	}

	// 本月活跃会员数（有过消费的会员）
	var activeMembers int64
	if err := h.db.Model(&models.Order{}).
		Where("status = ? AND created_at >= ?", "completed", firstDayOfMonth).
		Distinct("member_id").
		Count(&activeMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count active members", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"monthlyRevenue":    monthlyRevenue,
		"monthlyNewMembers": monthlyNewMembers,
		"monthlyOrders":     monthlyOrders,
		"activeMembers":     activeMembers,
	}, ""))
}

// GetProductSalesOverview returns product sales overview
// GET /api/dashboard/product-sales
func (h *DashboardHandler) GetProductSalesOverview(c *gin.Context) {
	// 获取时间范围参数，默认近30天
	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
			if days > 365 {
				days = 365 // 限制最多一年
			}
		}
	}

	now := time.Now()
	startDate := now.AddDate(0, 0, -days)

	type ProductSales struct {
		ProductID    uint    `json:"product_id"`
		ProductName  string  `json:"product_name"`
		SalesCount   int64   `json:"sales_count"`
		TotalRevenue float64 `json:"total_revenue"`
	}

	var topProducts []ProductSales

	// 统计热销商品
	if err := h.db.Model(&models.Order{}).
		Select("physical_products.id as product_id, physical_products.name as product_name, COUNT(orders.id) as sales_count, COALESCE(SUM(orders.actual_paid), 0) as total_revenue").
		Joins("JOIN physical_products ON physical_products.id = orders.product_id").
		Where("orders.type = ? AND orders.status = ? AND orders.created_at >= ?", "product", "completed", startDate).
		Group("physical_products.id, physical_products.name").
		Order("sales_count DESC").
		Limit(5).
		Scan(&topProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to get product sales", err.Error()))
		return
	}

	// 统计总销售额和总销量
	var totalRevenue float64
	var totalSales int64
	if err := h.db.Model(&models.Order{}).
		Where("type = ? AND status = ? AND created_at >= ?", "product", "completed", startDate).
		Select("COALESCE(SUM(actual_paid), 0)").
		Scan(&totalRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to calculate total revenue", err.Error()))
		return
	}

	if err := h.db.Model(&models.Order{}).
		Where("type = ? AND status = ? AND created_at >= ?", "product", "completed", startDate).
		Count(&totalSales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count total sales", err.Error()))
		return
	}

	// 统计库存预警（库存低于10的商品数）
	var lowStockCount int64
	if err := h.db.Model(&models.PhysicalProduct{}).
		Where("stock < ? AND is_active = ?", 10, true).
		Count(&lowStockCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to count low stock", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"topProducts":   topProducts,
		"totalRevenue":  totalRevenue,
		"totalSales":    totalSales,
		"lowStockCount": lowStockCount,
		"periodDays":    days,
	}, ""))
}
