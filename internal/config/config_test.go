package config

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestLoadConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret")
	os.Setenv("REDIS_HOST", "test-redis-host")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("DB_HOST", "test-db-host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("ENVIRONMENT", "test")

	// Clean up environment variables after test
	defer func() {
		os.Unsetenv("JWT_ACCESS_SECRET")
		os.Unsetenv("JWT_REFRESH_SECRET")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("ENVIRONMENT")
	}()

	// Create a test logger
	logger := zap.NewNop()

	// Load configuration
	cfg := LoadConfig(logger)

	// Test JWT configuration
	if cfg.JWT.AccessSecret != "test-access-secret" {
		t.Errorf("Expected JWT.AccessSecret to be 'test-access-secret', got '%s'", cfg.JWT.AccessSecret)
	}
	if cfg.JWT.RefreshSecret != "test-refresh-secret" {
		t.Errorf("Expected JWT.RefreshSecret to be 'test-refresh-secret', got '%s'", cfg.JWT.RefreshSecret)
	}

	// Test Redis configuration
	if cfg.Redis.Host != "test-redis-host" {
		t.Errorf("Expected Redis.Host to be 'test-redis-host', got '%s'", cfg.Redis.Host)
	}
	if cfg.Redis.Port != 6380 {
		t.Errorf("Expected Redis.Port to be 6380, got %d", cfg.Redis.Port)
	}

	// Test Database configuration
	if cfg.DB.Host != "test-db-host" {
		t.Errorf("Expected DB.Host to be 'test-db-host', got '%s'", cfg.DB.Host)
	}
	if cfg.DB.Port != "5433" {
		t.Errorf("Expected DB.Port to be '5433', got '%s'", cfg.DB.Port)
	}
	if cfg.DB.Name != "test_db" {
		t.Errorf("Expected DB.Name to be 'test_db', got '%s'", cfg.DB.Name)
	}
	if cfg.DB.User != "test_user" {
		t.Errorf("Expected DB.User to be 'test_user', got '%s'", cfg.DB.User)
	}
	if cfg.DB.Password != "test_password" {
		t.Errorf("Expected DB.Password to be 'test_password', got '%s'", cfg.DB.Password)
	}

	// Test Server configuration
	if cfg.Server.Port != "9090" {
		t.Errorf("Expected Server.Port to be '9090', got '%s'", cfg.Server.Port)
	}

	// Test Environment
	if cfg.Environment != "test" {
		t.Errorf("Expected Environment to be 'test', got '%s'", cfg.Environment)
	}
}

func TestConfigHelperMethods(t *testing.T) {
	// Set up test environment with required JWT secrets
	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret")
	os.Setenv("ENVIRONMENT", "development")
	defer func() {
		os.Unsetenv("JWT_ACCESS_SECRET")
		os.Unsetenv("JWT_REFRESH_SECRET")
		os.Unsetenv("ENVIRONMENT")
	}()

	logger := zap.NewNop()
	cfg := LoadConfig(logger)

	// Test helper methods
	if !cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return true for development environment")
	}
	if cfg.IsProduction() {
		t.Error("Expected IsProduction() to return false for development environment")
	}

	// Test production environment
	os.Setenv("ENVIRONMENT", "production")
	cfg = LoadConfig(logger)

	if cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return false for production environment")
	}
	if !cfg.IsProduction() {
		t.Error("Expected IsProduction() to return true for production environment")
	}
}

func TestConfigValidation(t *testing.T) {
	// Test validation logic directly
	cfg := &Config{}
	
	// Test with empty JWT secrets
	if err := validateConfig(cfg); err == nil {
		t.Error("Expected validation to fail with empty JWT secrets")
	}
	
	// Test with default JWT secrets
	cfg.JWT.AccessSecret = "your-access-secret-key"
	cfg.JWT.RefreshSecret = "your-refresh-secret-key"
	if err := validateConfig(cfg); err == nil {
		t.Error("Expected validation to fail with default JWT secrets")
	}
	
	// Test with valid JWT secrets but missing other fields
	cfg.JWT.AccessSecret = "valid-secret"
	cfg.JWT.RefreshSecret = "valid-secret"
	if err := validateConfig(cfg); err == nil {
		t.Error("Expected validation to fail with missing DB fields")
	}
	
	// Test with all required fields
	cfg.DB.Host = "localhost"
	cfg.DB.Port = "5432"
	cfg.DB.Name = "testdb"
	cfg.DB.User = "testuser"
	cfg.Server.Port = "8080"
	
	if err := validateConfig(cfg); err != nil {
		t.Errorf("Expected validation to pass with all required fields, got error: %v", err)
	}
}

func TestEnvironmentVariablePrecedence(t *testing.T) {
	// Test that environment variables take precedence over defaults
	os.Setenv("JWT_ACCESS_SECRET", "env-secret")
	os.Setenv("JWT_REFRESH_SECRET", "env-refresh-secret")
	os.Setenv("REDIS_HOST", "env-redis-host")
	os.Setenv("DB_HOST", "env-db-host")
	os.Setenv("SERVER_PORT", "9999")
	defer func() {
		os.Unsetenv("JWT_ACCESS_SECRET")
		os.Unsetenv("JWT_REFRESH_SECRET")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	logger := zap.NewNop()
	cfg := LoadConfig(logger)

	// Verify environment variables override defaults
	if cfg.JWT.AccessSecret != "env-secret" {
		t.Errorf("Expected JWT.AccessSecret to be 'env-secret', got '%s'", cfg.JWT.AccessSecret)
	}
	if cfg.JWT.RefreshSecret != "env-refresh-secret" {
		t.Errorf("Expected JWT.RefreshSecret to be 'env-refresh-secret', got '%s'", cfg.JWT.RefreshSecret)
	}
	if cfg.Redis.Host != "env-redis-host" {
		t.Errorf("Expected Redis.Host to be 'env-redis-host', got '%s'", cfg.Redis.Host)
	}
	if cfg.DB.Host != "env-db-host" {
		t.Errorf("Expected DB.Host to be 'env-db-host', got '%s'", cfg.DB.Host)
	}
	if cfg.Server.Port != "9999" {
		t.Errorf("Expected Server.Port to be '9999', got '%s'", cfg.Server.Port)
	}

	// Test debug method
	debug := cfg.DebugConfigSources()
	if !debug["jwt_access_secret_set"].(bool) {
		t.Error("Expected jwt_access_secret_set to be true")
	}
	if !debug["redis_host_set"].(bool) {
		t.Error("Expected redis_host_set to be true")
	}
}
