package services

import (
	"os"

	"chefly/models"

	"github.com/rs/zerolog"
)

// AuditLogger handles structured audit logging
type AuditLogger struct {
	logger  zerolog.Logger
	enabled bool
}

// NewAuditLogger creates a new audit logger instance
func NewAuditLogger(enabled bool, level, format string) *AuditLogger {
	// Determine output format
	var logger zerolog.Logger
	if format == "pretty" {
		// Pretty format for local development (colored, human-readable)
		output := zerolog.ConsoleWriter{Out: os.Stdout}
		logger = zerolog.New(output).With().Timestamp().Logger()
	} else {
		// JSON format for production (structured, machine-readable)
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Set log level
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return &AuditLogger{
		logger:  logger,
		enabled: enabled,
	}
}

// Info logs an informational audit event
func (a *AuditLogger) Info(eventType, message string, context *models.AuditContext) {
	if !a.enabled {
		return
	}

	event := a.logger.Info().
		Str("event_type", eventType).
		Str("message", message)

	if context != nil {
		event = a.addContext(event, context)
	}

	event.Send()
}

// Warn logs a warning audit event
func (a *AuditLogger) Warn(eventType, message string, context *models.AuditContext) {
	if !a.enabled {
		return
	}

	event := a.logger.Warn().
		Str("event_type", eventType).
		Str("message", message)

	if context != nil {
		event = a.addContext(event, context)
	}

	event.Send()
}

// Error logs an error audit event
func (a *AuditLogger) Error(eventType, message string, err error, context *models.AuditContext) {
	if !a.enabled {
		return
	}

	event := a.logger.Error().
		Str("event_type", eventType).
		Str("message", message)

	if err != nil {
		event = event.Err(err)
	}

	if context != nil {
		event = a.addContext(event, context)
	}

	event.Send()
}

// Debug logs a debug audit event
func (a *AuditLogger) Debug(eventType, message string, context *models.AuditContext) {
	if !a.enabled {
		return
	}

	event := a.logger.Debug().
		Str("event_type", eventType).
		Str("message", message)

	if context != nil {
		event = a.addContext(event, context)
	}

	event.Send()
}

// Fatal logs a fatal audit event and exits
func (a *AuditLogger) Fatal(eventType, message string, err error, context *models.AuditContext) {
	if !a.enabled {
		return
	}

	event := a.logger.Fatal().
		Str("event_type", eventType).
		Str("message", message)

	if err != nil {
		event = event.Err(err)
	}

	if context != nil {
		event = a.addContext(event, context)
	}

	event.Send()
}

// addContext adds audit context fields to a log event
func (a *AuditLogger) addContext(event *zerolog.Event, context *models.AuditContext) *zerolog.Event {
	if context.UserID != "" {
		event = event.Str("user_id", context.UserID)
	}
	if context.Username != "" {
		event = event.Str("username", context.Username)
	}
	if context.Email != "" {
		event = event.Str("email", context.Email)
	}
	if context.EmailAttempted != "" {
		event = event.Str("email_attempted", context.EmailAttempted)
	}
	if context.IsAdmin {
		event = event.Bool("is_admin", context.IsAdmin)
	}
	if context.IPAddress != "" {
		event = event.Str("ip_address", context.IPAddress)
	}
	if context.UserAgent != "" {
		event = event.Str("user_agent", context.UserAgent)
	}
	if context.RequestID != "" {
		event = event.Str("request_id", context.RequestID)
	}
	if context.Endpoint != "" {
		event = event.Str("endpoint", context.Endpoint)
	}
	if context.Method != "" {
		event = event.Str("method", context.Method)
	}
	if context.StatusCode != 0 {
		event = event.Int("status_code", context.StatusCode)
	}
	if context.DurationMs != 0 {
		event = event.Int64("duration_ms", context.DurationMs)
	}
	if context.Metadata != nil && len(context.Metadata) > 0 {
		event = event.Interface("metadata", context.Metadata)
	}

	return event
}
