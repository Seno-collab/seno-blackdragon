package config

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	JwtAccessSecret  string `mapstructure:"jwt_access_secret"`
	JwtRefreshSecret string `mapstructure:"jwt_refresh_secret"`
	RedisHost     string `mapstructure:"redis_host"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisPort     int    `mapstructure:"redis_port"`
	RedisDB       int    `mapstructure:"redis_db"`
	DBName     string `mapstructure:"db_name"`
	DBHost     string `mapstructure:"db_host"`
	DBPort string `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBSSLMode  string `mapstructure:"db_sslmode"`
	ServerPort string `mapstructure:"server_port"`
	ServerHost string `mapstructure:"server_host"`
	Environment string `mapstructure:"environment"`
}

func LoadConfig(logger *zap.Logger) *Config {
	// Set default values
	setDefaults()

	// Set environment variable key replacer to convert dots to underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// Enable automatic environment variable reading with highest priority
	viper.AutomaticEnv()
	
	// Ensure environment variables take precedence over config files
	viper.SetEnvPrefix("")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Warn("Error reading .env file", zap.Error(err))
		}
	} else {
		logger.Info(".env file loaded for backward compatibility")
	}

	// Create config struct
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		logger.Fatal("Unable to decode config", zap.Error(err))
	}

	logger.Info("Configuration loaded successfully", 
		zap.String("environment", cfg.Environment),
		zap.String("server_port", cfg.ServerPort),
		zap.String("db_host", cfg.DBHost))
	return cfg
}

func setDefaults() {
	// JWT defaults
	viper.SetDefault("jwt_access_secret", "your-access-secret-key")
	viper.SetDefault("jwt_refresh_secret", "your-refresh-secret-key")

	// Redis defaults
	viper.SetDefault("redis_host", "localhost")
	viper.SetDefault("redis_port", 6379)
	viper.SetDefault("redis_db", 0)
	viper.SetDefault("redis_password", "")

	// Database defaults
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", "5432")
	viper.SetDefault("db_name", "seno_blackdragon")
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "")
	viper.SetDefault("db_sslmode", "disable")

	// Server defaults
	viper.SetDefault("server_host", "0.0.0.0")
	viper.SetDefault("server_port", "8080")

	// Environment
	viper.SetDefault("environment", "development")
}

// GetString returns a string value from config
func (c *Config) GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns an int value from config
func (c *Config) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a bool value from config
func (c *Config) GetBool(key string) bool {
	return viper.GetBool(key)
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
