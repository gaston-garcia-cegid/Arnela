package middleware

import (
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate response time
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Log the request
		fields := map[string]interface{}{
			"method":     method,
			"path":       path,
			"status":     statusCode,
			"duration":   duration.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// Add user ID if authenticated
		if userID, exists := c.Get("userID"); exists {
			fields["user_id"] = userID
		}

		// Log based on status code
		if statusCode >= 500 {
			log.Error("Server error", nil, fields)
		} else if statusCode >= 400 {
			log.Warn("Client error", fields)
		} else {
			log.Info("Request processed", fields)
		}
	}
}
