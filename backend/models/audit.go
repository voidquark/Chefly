package models

import "time"

// AuditContext contains contextual information for audit logs
type AuditContext struct {
	UserID         string                 `json:"user_id,omitempty"`
	Username       string                 `json:"username,omitempty"`
	Email          string                 `json:"email,omitempty"`
	EmailAttempted string                 `json:"email_attempted,omitempty"` // For failed auth attempts
	IsAdmin        bool                   `json:"is_admin,omitempty"`
	IPAddress      string                 `json:"ip_address,omitempty"`
	UserAgent      string                 `json:"user_agent,omitempty"`
	RequestID      string                 `json:"request_id"`
	Endpoint       string                 `json:"endpoint,omitempty"`
	Method         string                 `json:"method,omitempty"`
	StatusCode     int                    `json:"status_code,omitempty"`
	DurationMs     int64                  `json:"duration_ms,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// AuditLog represents a complete audit log entry
type AuditLog struct {
	Timestamp time.Time     `json:"timestamp"`
	Level     string        `json:"level"`
	EventType string        `json:"event_type"`
	Message   string        `json:"message"`
	Error     string        `json:"error,omitempty"`
	Context   *AuditContext `json:"context,omitempty"`
}
