package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"
const requestIDContextKey = "request_id"

// RequestID ensures every response has X-Request-ID (echo client value or generate UUID).
// Useful for correlating logs with frontend traces (e.g. Sentry) when the client forwards the header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}
		c.Writer.Header().Set(requestIDHeader, id)
		c.Set(requestIDContextKey, id)
		c.Next()
	}
}
