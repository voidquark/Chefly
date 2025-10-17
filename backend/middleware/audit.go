package middleware

import (
	"time"

	"chefly/models"
	"chefly/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditMiddleware logs all API requests/responses
func AuditMiddleware(auditLogger *services.AuditLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()

		// Store in context for handlers to use
		c.Set("request_id", requestID)
		c.Set("audit_logger", auditLogger)

		// Process request
		c.Next()

		// Log after response
		duration := time.Since(start).Milliseconds()

		context := &models.AuditContext{
			RequestID:  requestID,
			Endpoint:   c.Request.URL.Path,
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			DurationMs: duration,
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		// Add user info if authenticated
		if userID, exists := c.Get("user_id"); exists {
			context.UserID = userID.(string)
		}
		if email, exists := c.Get("email"); exists {
			context.Email = email.(string)
		}
		if username, exists := c.Get("username"); exists {
			context.Username = username.(string)
		}
		if isAdmin, exists := c.Get("is_admin"); exists {
			context.IsAdmin = isAdmin.(bool)
		}

		auditLogger.Info("api.request", "API request processed", context)
	}
}
