package main

import (
	"log"
	"time"

	"github.com/yourusername/Task_Management/internal/api"
	"github.com/yourusername/Task_Management/internal/db"
	"github.com/yourusername/Task_Management/internal/config"
)

func main() {
	// ⚠️ Hardcoded config (not recommended for production)
	cfg := &config.Config{
		ServerAddress:   ":8080",
		DatabaseURL:     "postgres://postgres:postgres@localhost:5432/task_db?sslmode=disable",
		JWTSecret:       "supersecretkey123",
		TokenExpiration: 24 * time.Hour,
	}
	cfg.RateLimit.Period = time.Minute
	cfg.RateLimit.Limit = 60

	// Initialize database
	database, err := db.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Start API server
	router := api.SetupRouter(cfg, database.DB)
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
