package main

import (
	"log"
	"net/http"

	"smartspa-admin/internal/db"
	"smartspa-admin/internal/handlers"
	"smartspa-admin/internal/response"

	"github.com/gin-gonic/gin"
)

func main() {
	if _, err := db.Init("spa_management.db"); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

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

	// Dashboard APIs
	r.GET("/api/dashboard/stats", handlers.GetDashboardStats)
	r.GET("/api/appointments", handlers.ListAppointments)
	r.GET("/api/fission/ranking", handlers.GetFissionRanking)
	r.POST("/api/appointments", handlers.CreateAppointment)
	r.PUT("/api/appointments/:id/cancel", handlers.CancelAppointment)
	r.PUT("/api/appointments/:id/complete", handlers.CompleteAppointment)

	// Technician APIs
	r.GET("/api/technicians", handlers.ListTechnicians)
	r.POST("/api/technicians", handlers.CreateTechnician)
	r.PUT("/api/technicians/:id", handlers.UpdateTechnician)
	r.DELETE("/api/technicians/:id", handlers.DeleteTechnician)

	// Schedule APIs
	r.GET("/api/schedules", handlers.GetSchedules)
	r.POST("/api/schedules/batch", handlers.BatchSetSchedule)

	// Service Item APIs
	r.GET("/api/services", handlers.ListServiceItems)
	r.POST("/api/services", handlers.CreateServiceItem)
	r.PUT("/api/services/:id", handlers.UpdateServiceItem)
	r.DELETE("/api/services/:id", handlers.DeleteServiceItem)

	// Member APIs
	r.GET("/api/members", handlers.ListMembers)
	r.POST("/api/members", handlers.CreateMember)

	// AI APIs
	r.GET("/api/ai/report", handlers.GenerateAIReport)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
