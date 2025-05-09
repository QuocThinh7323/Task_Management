package middleware

import (
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/yourusername/Task_Management/internal/config"
)

// RateLimitMiddleware implements rate limiting for API requests
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Create a rate limiter with the specified limit and period
	rate := limiter.Rate{
		Period: cfg.RateLimit.Period,
		Limit:  cfg.RateLimit.Limit,
	}
	
	// Create a memory store with default options
	store := memory.NewStore()
	
	// Create a new limiter instance
	rateLimiter := limiter.New(store, rate)
	
	return func(c *gin.Context) {
		// Use the client IP as the identifier for rate limiting
		clientIP := c.ClientIP()
		
		// Check if the request exceeds the rate limit
		context, err := rateLimiter.Get(c, clientIP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking rate limit"})
			c.Abort()
			return
		}
		
		// Set headers to inform client of rate limit status
		c.Header("X-RateLimit-Limit", string(context.Limit))
		c.Header("X-RateLimit-Remaining", string(context.Remaining))
		c.Header("X-RateLimit-Reset", string(context.Reset))
		
		// If rate limit exceeded, return 429 Too Many Requests
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
				"limit": context.Limit,
				"reset": context.Reset - time.Now().Unix(),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
