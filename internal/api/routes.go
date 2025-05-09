package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/yourusername/Task_Management/internal/api/handlers"
	"github.com/yourusername/Task_Management/internal/api/middleware"
	"github.com/yourusername/Task_Management/internal/config"
	"github.com/yourusername/Task_Management/internal/models"
)

// SetupRouter configures the API routes
func SetupRouter(cfg *config.Config, db *sqlx.DB) *gin.Engine {
	// Create a new Gin router
	router := gin.New()
	
	// Set up logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	// Use middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.RateLimitMiddleware(cfg))
	
	// Create repositories
	userRepo := models.NewUserRepository(db)
	taskRepo := models.NewTaskRepository(db)
	categoryRepo := models.NewCategoryRepository(db)
	
	// Create handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg)
	userHandler := handlers.NewUserHandler(userRepo)
	taskHandler := handlers.NewTaskHandler(taskRepo, categoryRepo)
	
	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	api.Use(middleware.AuditLogger(db))
	
	// User routes
	api.GET("/users", middleware.RequireRole("admin"), userHandler.GetUsers)
	api.GET("/users/:id", userHandler.GetUser)
	
	// Task routes
	api.POST("/tasks", taskHandler.CreateTask)
	api.GET("/tasks", taskHandler.GetTasks)
	api.GET("/tasks/:id", taskHandler.GetTask)
	api.PUT("/tasks/:id", taskHandler.UpdateTask)
	api.DELETE("/tasks/:id", taskHandler.DeleteTask)
	
	// Category routes
	api.GET("/categories", taskHandler.GetCategories)
	api.POST("/categories", middleware.RequireRole("admin"), taskHandler.CreateCategory)
	api.DELETE("/categories/:id", middleware.RequireRole("admin"), taskHandler.DeleteCategory)
	
	return router
}