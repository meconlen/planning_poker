package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	Port            string        `json:"port"`
	Host            string        `json:"host"`
	ReadTimeout     time.Duration `json:"readTimeout"`
	WriteTimeout    time.Duration `json:"writeTimeout"`
	IdleTimeout     time.Duration `json:"idleTimeout"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`

	// WebSocket configuration
	AllowedOrigins []string `json:"allowedOrigins"`
	MaxMessageSize int64    `json:"maxMessageSize"`

	// Session configuration
	SessionTimeout     time.Duration `json:"sessionTimeout"`
	MaxSessionsPerUser int           `json:"maxSessionsPerUser"`

	// Logging configuration
	LogLevel  string `json:"logLevel"`
	LogFormat string `json:"logFormat"`

	// Development settings
	IsDevelopment bool `json:"isDevelopment"`
	EnablePprof   bool `json:"enablePprof"`
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	config := &Config{
		// Default values
		Port:               getEnv("PORT", "8080"),
		Host:               getEnv("HOST", ""),
		ReadTimeout:        getDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:       getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:        getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:    getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		AllowedOrigins:     getStringSliceEnv("ALLOWED_ORIGINS", []string{"*"}),
		MaxMessageSize:     getInt64Env("MAX_MESSAGE_SIZE", 1024),
		SessionTimeout:     getDurationEnv("SESSION_TIMEOUT", 24*time.Hour),
		MaxSessionsPerUser: getIntEnv("MAX_SESSIONS_PER_USER", 10),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFormat:          getEnv("LOG_FORMAT", "text"),
		IsDevelopment:      getBoolEnv("DEVELOPMENT", false),
		EnablePprof:        getBoolEnv("ENABLE_PPROF", false),
	}

	// In development mode, be more permissive
	if config.IsDevelopment {
		config.AllowedOrigins = []string{"*"}
		config.LogLevel = "debug"
		config.EnablePprof = true
	}

	return config
}

// Address returns the server address (host:port)
func (c *Config) Address() string {
	return c.Host + ":" + c.Port
}

// IsProductionMode returns true if not in development mode
func (c *Config) IsProductionMode() bool {
	return !c.IsDevelopment
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
