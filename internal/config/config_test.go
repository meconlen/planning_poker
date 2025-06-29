package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	config := Load()

	// Test default values
	if config.Port != "8080" {
		t.Errorf("Expected default port '8080', got %s", config.Port)
	}

	if config.Host != "" {
		t.Errorf("Expected default host '', got %s", config.Host)
	}

	if config.ReadTimeout != 15*time.Second {
		t.Errorf("Expected default read timeout 15s, got %v", config.ReadTimeout)
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected default log level 'info', got %s", config.LogLevel)
	}

	if config.IsDevelopment != false {
		t.Errorf("Expected default development mode false, got %v", config.IsDevelopment)
	}

	if len(config.AllowedOrigins) != 1 || config.AllowedOrigins[0] != "*" {
		t.Errorf("Expected default allowed origins ['*'], got %v", config.AllowedOrigins)
	}
}

func TestLoad_EnvironmentOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9000")
	os.Setenv("HOST", "localhost")
	os.Setenv("READ_TIMEOUT", "30s")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DEVELOPMENT", "false") // Set to false to test environment overrides
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,https://example.com")
	os.Setenv("MAX_MESSAGE_SIZE", "2048")

	defer func() {
		os.Clearenv()
	}()

	config := Load()

	if config.Port != "9000" {
		t.Errorf("Expected port '9000', got %s", config.Port)
	}

	if config.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got %s", config.Host)
	}

	if config.ReadTimeout != 30*time.Second {
		t.Errorf("Expected read timeout 30s, got %v", config.ReadTimeout)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got %s", config.LogLevel)
	}

	if config.IsDevelopment != false {
		t.Errorf("Expected development mode false, got %v", config.IsDevelopment)
	}

	expectedOrigins := []string{"http://localhost:3000", "https://example.com"}
	if len(config.AllowedOrigins) != 2 || config.AllowedOrigins[0] != expectedOrigins[0] || config.AllowedOrigins[1] != expectedOrigins[1] {
		t.Errorf("Expected allowed origins %v, got %v", expectedOrigins, config.AllowedOrigins)
	}

	if config.MaxMessageSize != 2048 {
		t.Errorf("Expected max message size 2048, got %d", config.MaxMessageSize)
	}
}

func TestLoad_DevelopmentMode(t *testing.T) {
	os.Setenv("DEVELOPMENT", "true")
	defer os.Clearenv()

	config := Load()

	if !config.IsDevelopment {
		t.Error("Expected development mode to be true")
	}

	// In development mode, allowed origins should be permissive
	if len(config.AllowedOrigins) != 1 || config.AllowedOrigins[0] != "*" {
		t.Errorf("Expected development mode to set allowed origins to ['*'], got %v", config.AllowedOrigins)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Expected development mode to set log level to 'debug', got %s", config.LogLevel)
	}

	if !config.EnablePprof {
		t.Error("Expected development mode to enable pprof")
	}
}

func TestAddress(t *testing.T) {
	config := &Config{
		Host: "localhost",
		Port: "8080",
	}

	expected := "localhost:8080"
	if address := config.Address(); address != expected {
		t.Errorf("Expected address '%s', got '%s'", expected, address)
	}

	// Test with empty host
	config.Host = ""
	expected = ":8080"
	if address := config.Address(); address != expected {
		t.Errorf("Expected address '%s', got '%s'", expected, address)
	}
}

func TestIsProductionMode(t *testing.T) {
	config := &Config{IsDevelopment: false}
	if !config.IsProductionMode() {
		t.Error("Expected production mode to be true when development is false")
	}

	config.IsDevelopment = true
	if config.IsProductionMode() {
		t.Error("Expected production mode to be false when development is true")
	}
}

func TestGetEnvHelpers(t *testing.T) {
	// Test getEnv
	os.Setenv("TEST_STRING", "test_value")
	if value := getEnv("TEST_STRING", "default"); value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}

	if value := getEnv("NONEXISTENT", "default"); value != "default" {
		t.Errorf("Expected 'default', got '%s'", value)
	}

	// Test getIntEnv
	os.Setenv("TEST_INT", "42")
	if value := getIntEnv("TEST_INT", 0); value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	if value := getIntEnv("NONEXISTENT_INT", 10); value != 10 {
		t.Errorf("Expected 10, got %d", value)
	}

	// Test getBoolEnv
	os.Setenv("TEST_BOOL", "true")
	if value := getBoolEnv("TEST_BOOL", false); value != true {
		t.Errorf("Expected true, got %v", value)
	}

	if value := getBoolEnv("NONEXISTENT_BOOL", false); value != false {
		t.Errorf("Expected false, got %v", value)
	}

	// Test getDurationEnv
	os.Setenv("TEST_DURATION", "5m")
	if value := getDurationEnv("TEST_DURATION", time.Second); value != 5*time.Minute {
		t.Errorf("Expected 5m, got %v", value)
	}

	if value := getDurationEnv("NONEXISTENT_DURATION", time.Second); value != time.Second {
		t.Errorf("Expected 1s, got %v", value)
	}

	// Test getStringSliceEnv
	os.Setenv("TEST_SLICE", "a,b,c")
	expected := []string{"a", "b", "c"}
	if value := getStringSliceEnv("TEST_SLICE", []string{"default"}); !equalStringSlices(value, expected) {
		t.Errorf("Expected %v, got %v", expected, value)
	}

	defaultSlice := []string{"default"}
	if value := getStringSliceEnv("NONEXISTENT_SLICE", defaultSlice); !equalStringSlices(value, defaultSlice) {
		t.Errorf("Expected %v, got %v", defaultSlice, value)
	}

	os.Clearenv()
}

func TestInvalidEnvironmentValues(t *testing.T) {
	// Test invalid int
	os.Setenv("INVALID_INT", "not_a_number")
	if value := getIntEnv("INVALID_INT", 42); value != 42 {
		t.Errorf("Expected default value 42 for invalid int, got %d", value)
	}

	// Test invalid bool
	os.Setenv("INVALID_BOOL", "not_a_bool")
	if value := getBoolEnv("INVALID_BOOL", true); value != true {
		t.Errorf("Expected default value true for invalid bool, got %v", value)
	}

	// Test invalid duration
	os.Setenv("INVALID_DURATION", "not_a_duration")
	if value := getDurationEnv("INVALID_DURATION", time.Minute); value != time.Minute {
		t.Errorf("Expected default value 1m for invalid duration, got %v", value)
	}

	os.Clearenv()
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
