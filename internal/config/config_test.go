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
	if cfg.JwtAccessSecret != "test-access-secret" {
		t.Errorf("Expected JwtAccessSecret to be 'test-access-secret', got '%s'", cfg.JwtAccessSecret)
	}
	if cfg.JwtRefreshSecret != "test-refresh-secret" {
		t.Errorf("Expected JwtRefreshSecret to be 'test-refresh-secret', got '%s'", cfg.JwtRefreshSecret)
	}

	// Test Redis configuration
	if cfg.RedisHost != "test-redis-host" {
		t.Errorf("Expected RedisHost to be 'test-redis-host', got '%s'", cfg.RedisHost)
	}
	if cfg.RedisPort != 6380 {
		t.Errorf("Expected RedisPort to be 6380, got %d", cfg.RedisPort)
	}

	// Test Database configuration
	if cfg.DBHost != "test-db-host" {
		t.Errorf("Expected DBHost to be 'test-db-host', got '%s'", cfg.DBHost)
	}
	if cfg.DBPort != "5433" {
		t.Errorf("Expected DBPort to be '5433', got '%s'", cfg.DBPort)
	}
	if cfg.DBName != "test_db" {
		t.Errorf("Expected DBName to be 'test_db', got '%s'", cfg.DBName)
	}
	if cfg.DBUser != "test_user" {
		t.Errorf("Expected DBUser to be 'test_user', got '%s'", cfg.DBUser)
	}
	if cfg.DBPassword != "test_password" {
		t.Errorf("Expected DBPassword to be 'test_password', got '%s'", cfg.DBPassword)
	}

	// Test Server configuration
	if cfg.ServerPort != "9090" {
		t.Errorf("Expected ServerPort to be '9090', got '%s'", cfg.ServerPort)
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
	// Test that required fields are properly loaded
	cfg := &Config{}
	
	// Test with empty config (should use defaults or empty values)
	if cfg.JwtAccessSecret != "" {
		t.Errorf("Expected JwtAccessSecret to be empty, got '%s'", cfg.JwtAccessSecret)
	}
	
	// Test with environment variables set
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret")
	os.Setenv("DB_HOST", "test-host")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("SERVER_PORT", "8080")
	defer func() {
		os.Unsetenv("JWT_ACCESS_SECRET")
		os.Unsetenv("JWT_REFRESH_SECRET")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("SERVER_PORT")
	}()

	logger := zap.NewNop()
	cfg = LoadConfig(logger)
	
	// Verify that environment variables were loaded
	if cfg.JwtAccessSecret != "test-secret" {
		t.Errorf("Expected JwtAccessSecret to be 'test-secret', got '%s'", cfg.JwtAccessSecret)
	}
	if cfg.JwtRefreshSecret != "test-refresh-secret" {
		t.Errorf("Expected JwtRefreshSecret to be 'test-refresh-secret', got '%s'", cfg.JwtRefreshSecret)
	}
	if cfg.DBHost != "test-host" {
		t.Errorf("Expected DBHost to be 'test-host', got '%s'", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("Expected DBPort to be '5432', got '%s'", cfg.DBPort)
	}
	if cfg.DBName != "testdb" {
		t.Errorf("Expected DBName to be 'testdb', got '%s'", cfg.DBName)
	}
	if cfg.DBUser != "testuser" {
		t.Errorf("Expected DBUser to be 'testuser', got '%s'", cfg.DBUser)
	}
	if cfg.ServerPort != "8080" {
		t.Errorf("Expected ServerPort to be '8080', got '%s'", cfg.ServerPort)
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
	if cfg.JwtAccessSecret != "env-secret" {
		t.Errorf("Expected JwtAccessSecret to be 'env-secret', got '%s'", cfg.JwtAccessSecret)
	}
	if cfg.JwtRefreshSecret != "env-refresh-secret" {
		t.Errorf("Expected JwtRefreshSecret to be 'env-refresh-secret', got '%s'", cfg.JwtRefreshSecret)
	}
	if cfg.RedisHost != "env-redis-host" {
		t.Errorf("Expected RedisHost to be 'env-redis-host', got '%s'", cfg.RedisHost)
	}
	if cfg.DBHost != "env-db-host" {
		t.Errorf("Expected DBHost to be 'env-db-host', got '%s'", cfg.DBHost)
	}
	if cfg.ServerPort != "9999" {
		t.Errorf("Expected ServerPort to be '9999', got '%s'", cfg.ServerPort)
	}

	// Test that environment variables were properly loaded
	if cfg.JwtAccessSecret == "" {
		t.Error("Expected JwtAccessSecret to be set")
	}
	if cfg.RedisHost == "" {
		t.Error("Expected RedisHost to be set")
	}
}
