package main

import (
	"log"
	"net/http"

	"server/internal/db"
	"server/internal/handlers"
	"server/internal/middleware"
	"server/internal/response"

	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.Init("spa_management.db")
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	// Initialize handlers
	dashboardHandler := handlers.NewDashboardHandler(database)

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Success(gin.H{"status": "ok"}, ""))
	})

	// Public routes (no authentication required)
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", handlers.Login)
	}

	// Protected routes (authentication required)
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		// Auth routes for authenticated users
		api.GET("/auth/me", handlers.GetCurrentUser)

		// Dashboard APIs (both manager and operator)
		api.GET("/dashboard/stats", handlers.GetDashboardStats)
		api.GET("/dashboard/revenue-trend", dashboardHandler.GetRevenueTrend)
		api.GET("/dashboard/service-ranking", dashboardHandler.GetServiceRanking)
		api.GET("/dashboard/product-sales", dashboardHandler.GetProductSalesOverview)
		api.GET("/dashboard/monthly-stats", dashboardHandler.GetMonthlyStats)
		api.GET("/dashboard/marketing", dashboardHandler.GetMarketingMetrics)

		// Appointments (both manager and operator)
		api.GET("/appointments", handlers.ListAppointments)
		api.POST("/appointments", handlers.CreateAppointment)
		api.PUT("/appointments/:id/cancel", handlers.CancelAppointment)
		api.PUT("/appointments/:id/complete", handlers.CompleteAppointment)

		// Fission ranking (both manager and operator)
		api.GET("/fission/ranking", dashboardHandler.GetFissionRanking)

		// Technicians (read for all, write for manager only)
		api.GET("/technicians", handlers.ListTechnicians)

		// Schedules (both manager and operator)
		api.GET("/schedules", handlers.GetSchedules)
		api.GET("/schedules/detail", handlers.GetTechnicianScheduleDetail)
		api.GET("/schedules/available-technicians", handlers.GetAvailableTechnicians)
		api.GET("/schedules/slots", handlers.GetTimeSlotsAvailability)
		api.POST("/schedules/batch", handlers.BatchSetSchedule)

		// Services (read for all, write for manager only)
		api.GET("/services", handlers.ListServiceItems)

		// Members (both manager and operator)
		api.GET("/members", handlers.ListMembers)
		api.POST("/members", handlers.CreateMember)

		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders", handlers.ListOrders)

		// Products (read for all, write for manager only)
		api.GET("/products", handlers.ListProducts)
		api.GET("/products/:id", handlers.GetProduct)
		api.GET("/products/stats", handlers.GetProductStats)

		// Inventory (both manager and operator can view and change)
		api.GET("/inventory/logs", handlers.ListInventoryLogs)
		api.GET("/inventory/products/:id/logs", handlers.GetProductInventoryLogs)
		api.POST("/inventory/change", handlers.CreateInventoryChange)
		api.POST("/inventory/batch-restock", handlers.BatchRestock)
		api.GET("/inventory/stats", handlers.GetInventoryStats)
	}

	// Manager-only routes
	managerAPI := r.Group("/api")
	managerAPI.Use(middleware.AuthRequired(), middleware.RequireManager())
	{
		// User management (manager only)
		managerAPI.POST("/auth/register", handlers.Register)
		managerAPI.GET("/auth/users", handlers.ListUsers)

		// Technician management (manager only)
		managerAPI.POST("/technicians", handlers.CreateTechnician)
		managerAPI.PUT("/technicians/:id", handlers.UpdateTechnician)
		managerAPI.DELETE("/technicians/:id", handlers.DeleteTechnician)

		// Service management (manager only)
		managerAPI.POST("/services", handlers.CreateServiceItem)
		managerAPI.PUT("/services/:id", handlers.UpdateServiceItem)
		managerAPI.DELETE("/services/:id", handlers.DeleteServiceItem)

		// Product management (manager only for create/update/delete)
		managerAPI.POST("/products", handlers.CreateProduct)
		managerAPI.PUT("/products/:id", handlers.UpdateProduct)
		managerAPI.DELETE("/products/:id", handlers.DeleteProduct)

		// Member AI profile (manager only)
		managerAPI.GET("/members/:id/ai-profile", handlers.GenerateMemberAIProfile)

		// AI report (manager only)
		managerAPI.GET("/ai/report", handlers.GenerateAIReport)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
