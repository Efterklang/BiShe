package middleware

import (
	"net/http"
	"strings"

	"server/internal/auth"
	"server/internal/response"

	"github.com/gin-gonic/gin"
)

// AuthRequired is a middleware that validates JWT token from Authorization header
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Authorization header required", nil))
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Invalid authorization header format", nil))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			if err == auth.ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token has expired", nil))
			} else {
				c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Invalid token", nil))
			}
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
