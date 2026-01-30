package middleware

import (
	"net/http"

	"server/internal/response"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that checks if the user has the required role
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by AuthRequired middleware)
		roleInterface, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Role information not found", nil))
			c.Abort()
			return
		}

		userRole, ok := roleInterface.(string)
		if !ok {
			c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Invalid role information", nil))
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles list
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		// User doesn't have required role
		c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Insufficient permissions", nil))
		c.Abort()
	}
}

// RequireManager is a convenience middleware that requires manager role
func RequireManager() gin.HandlerFunc {
	return RequireRole("manager")
}
