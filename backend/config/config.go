package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	DBPath                 string
	Port                   string
	JWTSecret              string
	ClaudeAPIKey           string
	ClaudeModel            string // Claude AI model (e.g: claude-3-haiku-20240307)
	OpenAIAPIKey           string
	OpenAIModel            string // OpenAI model (e.g: dall-e-3)
	ImageStoragePath       string
	Environment            string
	AuditLogEnabled        bool   // Enable audit logging
	AuditLogLevel          string // Log level: debug, info, warn, error
	AuditLogFormat         string // Log format: json or pretty
	RegistrationEnabled    bool
	RecipeGenerationLimit  string // Global recipe generation limit: "unlimited", "0", or number (e.g. "10")
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		DBPath:                getEnv("DB_PATH", "./data/chefly.db"),
		Port:                  getEnv("PORT", "8080"),
		JWTSecret:             getEnv("JWT_SECRET", "change-me-in-production"),
		ClaudeAPIKey:          getEnv("CLAUDE_API_KEY", ""),
		ClaudeModel:           getEnv("CLAUDE_MODEL", "claude-3-haiku-20240307"),
		OpenAIAPIKey:          getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:           getEnv("OPENAI_MODEL", "dall-e-3"),
		ImageStoragePath:      getEnv("IMAGE_STORAGE_PATH", "./data/images"),
		Environment:           getEnv("ENVIRONMENT", "development"),
		AuditLogEnabled:       getEnvBool("AUDIT_LOG_ENABLED", true),          // Default: enabled
		AuditLogLevel:         getEnv("AUDIT_LOG_LEVEL", "info"),              // Default: info
		AuditLogFormat:        getEnv("AUDIT_LOG_FORMAT", "json"),             // Default: json
		RegistrationEnabled:   getEnvBool("REGISTRATION_ENABLED", true),       // Default: enabled
		RecipeGenerationLimit: getEnv("RECIPE_GENERATION_LIMIT", "unlimited"), // Default: unlimited
	}
}

// getEnv gets environment variable with fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets boolean environment variable with fallback default
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}
