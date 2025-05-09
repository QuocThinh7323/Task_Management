package handlers

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/yourusername/Task_Management/internal/models"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userRepo *models.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *models.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetUsers returns all users (admin only)
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

// GetUser returns a specific user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get the requesting user's ID and role
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")
	
	// Only admins can view other users' details
	if userID.(int) != id && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}