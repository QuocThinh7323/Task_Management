package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress   string
	DatabaseURL     string
	JWTSecret       string
	TokenExpiration time.Duration
	RateLimit       struct {
		Period time.Duration
		Limit  int64
	}
}

func Load() (*Config, error) {
	// In a real application, you might use a library like viper
	// to load config from env vars, config files, etc.
	_ = godotenv.Load()

	
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = ":8080"
	}
	
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}
	
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}
	
	cfg := &Config{
		ServerAddress:   serverAddr,
		DatabaseURL:     dbURL,
		JWTSecret:       jwtSecret,
		TokenExpiration: 24 * time.Hour,
	}
	
	// Set rate limiting defaults
	cfg.RateLimit.Period = time.Minute
	cfg.RateLimit.Limit = 60 // 60 requests per minute
	
	return cfg, nil
}
