package handlers

import (
	
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yourusername/Task_Management/internal/config"
	"github.com/yourusername/Task_Management/internal/models"
	"github.com/yourusername/Task_Management/internal/utils"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userRepo *models.UserRepository
	config   *config.Config
	validate *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *models.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		config:   cfg,
		validate: validator.New(),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors.Error()})
		return
	}
	
// Check if username already exists
existingUser, err := h.userRepo.FindByUsername(req.Username)
if err != nil {
    // Đây là phần quan trọng: Nếu lỗi là "no rows in result set", 
    // có nghĩa là username chưa tồn tại, nên tiếp tục quá trình đăng ký
    if err.Error() == "sql: no rows in result set" {
        // Username chưa tồn tại, tiếp tục đăng ký
    } else {
        // Lỗi database thực sự
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        return
    }
} else if existingUser != nil {
    // Username đã tồn tại
    c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
    return
}
	
	// Create new user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "user", // Default role for new users
	}
	
	if err := h.userRepo.Create(user, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	
	// Generate JWT token
	token, err := utils.GenerateToken(user, h.config.JWTSecret, h.config.TokenExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Find user by username
	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	
	// Verify password
	if !h.userRepo.CheckPassword(user, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	
	// Generate JWT token
	token, err := utils.GenerateToken(user, h.config.JWTSecret, h.config.TokenExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
