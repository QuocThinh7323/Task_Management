package middleware

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// responseWriter wraps the gin ResponseWriter to capture the response
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware logs requests and responses
func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Read the request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		
		// Capture the response
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w
		
		// Process request
		c.Next()
		
		// Calculate request duration
		duration := time.Since(start)
		
		// Log the request details
		entry := logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration.String(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})
		
		// Add user info if available
		if userID, exists := c.Get("userID"); exists {
			entry = entry.WithField("user_id", userID)
		}
		
		// Log based on status code
		if c.Writer.Status() >= 500 {
			entry.Error("Server error")
		} else if c.Writer.Status() >= 400 {
			entry.Warn("Client error")
		} else {
			entry.Info("Request processed")
		}
	}
}

// AuditLogger records significant actions to the database
func AuditLogger(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request first
		c.Next()
		
		// Only log specific actions (POST, PUT, DELETE) for auditing
		if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "DELETE" {
			return
		}
		
		// Extract user ID if available
		var userID *int
		if id, exists := c.Get("userID"); exists {
			if val, ok := id.(int); ok {
				userID = &val
			}
		}
		
		// Determine the entity type and ID from the path
		pathParts := strings.Split(c.Request.URL.Path, "/")
		entityType := ""
		var entityID *int
		
		if len(pathParts) >= 2 {
			entityType = pathParts[1] // e.g., "tasks", "users"
			
			// Try to extract entity ID if present
			if len(pathParts) >= 3 {
				if id, err := strconv.Atoi(pathParts[2]); err == nil {
					entityID = &id
				}
			}
		}
		
		// Determine the action based on HTTP method
		action := ""
		switch c.Request.Method {
		case "POST":
			action = "create"
		case "PUT":
			action = "update"
		case "DELETE":
			action = "delete"
		}
		
		// Insert audit log
		_, err := db.Exec(
			`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, ip_address, created_at)
			 VALUES ($1, $2, $3, $4, $5, NOW())`,
			userID, action, entityType, entityID, c.ClientIP(),
		)
		
		if err != nil {
			logrus.WithError(err).Error("Failed to insert audit log")
		}
	}
}