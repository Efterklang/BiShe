package handlers

import (
	"net/http"

	"smartspa-admin/internal/auth"
	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"
	"smartspa-admin/internal/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=manager operator"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse represents user information without sensitive data
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// Register creates a new user account (manager only)
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid request: "+err.Error(), nil))
		return
	}

	// Validate password strength
	if !isValidPassword(req.Password) {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Password must be at least 8 characters and contain uppercase, lowercase, and numbers", nil))
		return
	}

	// Check if username already exists
	var existingUser models.User
	database := db.GetDB()
	if err := database.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, response.Error(http.StatusConflict, "Username already exists", nil))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to process password", nil))
		return
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		IsActive:     true,
	}

	if err := database.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to create user", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, "User created successfully"))
}

// Login authenticates a user and returns a JWT token
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid request: "+err.Error(), nil))
		return
	}

	// Find user by username
	var user models.User
	database := db.GetDB()
	if err := database.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Invalid username or password", nil))
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Account has been deactivated", nil))
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Invalid username or password", nil))
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to generate token", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(LoginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, "Login successful"))
}

// GetCurrentUser returns the current authenticated user's information
func GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by AuthRequired middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "User not authenticated", nil))
		return
	}

	// Fetch user from database
	var user models.User
	database := db.GetDB()
	if err := database.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "User not found", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, ""))
}

// ListUsers returns a list of all users (manager only)
func ListUsers(c *gin.Context) {
	var users []models.User
	database := db.GetDB()

	if err := database.Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch users", nil))
		return
	}

	// Convert to response format
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
			IsActive: user.IsActive,
		}
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"users": userResponses,
		"total": len(userResponses),
	}, ""))
}

// isValidPassword checks if password meets security requirements
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}
