package handlers

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yourusername/Task_Management/internal/models"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	taskRepo     *models.TaskRepository
	categoryRepo *models.CategoryRepository
	validate     *validator.Validate
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskRepo *models.TaskRepository, categoryRepo *models.CategoryRepository) *TaskHandler {
	return &TaskHandler{
		taskRepo:     taskRepo,
		categoryRepo: categoryRepo,
		validate:     validator.New(),
	}
}

// CreateTask handles task creation
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Validate the request
	if err := h.validate.Struct(task); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors.Error()})
		return
	}
	
	// Set user ID from authenticated user
	userID, _ := c.Get("userID")
	task.UserID = userID.(int)
	
	// Set default status if not provided
	if task.Status == "" {
		task.Status = "pending"
	}
	
	// Create the task
	if err := h.taskRepo.Create(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	
	c.JSON(http.StatusCreated, task)
}

// UpdateTask handles task updates
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	
	// Check if task exists
	existingTask, err := h.taskRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	
	// Get the requesting user's ID and role
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")
	
	// Only task owner or admin can update the task
	if existingTask.UserID != userID.(int) && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	
	// Bind the request body to update the task
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Preserve the ID and user ID
	updatedTask.ID = id
	updatedTask.UserID = existingTask.UserID
	
	// Update the task
	if err := h.taskRepo.Update(&updatedTask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	
	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask handles task deletion
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	
	// Check if task exists
	existingTask, err := h.taskRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	
	// Get the requesting user's ID and role
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")
	
	// Only task owner or admin can delete the task
	if existingTask.UserID != userID.(int) && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	
	// Delete the task
	if err := h.taskRepo.Delete(id, existingTask.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// GetTask returns a specific task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	
	// Get the task
	task, err := h.taskRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	
	// Get the requesting user's ID and role
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")
	
	// Only task owner or admin can view the task
	if task.UserID != userID.(int) && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	
	c.JSON(http.StatusOK, task)
}

// GetTasks returns all tasks for the authenticated user
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")
	
	var tasks []models.Task
	var err error
	
	// Admin can see all tasks, regular users only see their own
	if userRole.(string) == "admin" {
		tasks, err = h.taskRepo.ListAll()
	} else {
		tasks, err = h.taskRepo.ListByUser(userID.(int))
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	
	c.JSON(http.StatusOK, tasks)
}

// CreateCategory handles category creation (admin only)
func (h *TaskHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Validate the request
	if err := h.validate.Struct(category); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors.Error()})
		return
	}
	
	// Create the category
	if err := h.categoryRepo.Create(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	
	c.JSON(http.StatusCreated, category)
}

// GetCategories returns all categories
func (h *TaskHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	
	c.JSON(http.StatusOK, categories)
}

// DeleteCategory handles category deletion (admin only)
func (h *TaskHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	
	// Delete the category
	if err := h.categoryRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}