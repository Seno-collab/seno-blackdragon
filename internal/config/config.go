package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	JWT struct {
		AccessSecret  string
		RefreshSecret string
	}
	Redis struct {
		Host     string
		Password string
	}
	DB struct {
		Port     string
		Name     string
		Host     string
		User     string
		Password string
	}
	Server struct {
		Port string
	}
}

func LoadConfig(logger *zap.Logger) *Config {
	if err := godotenv.Load(); err != nil {
		logger.Error("Count not load .env file", zap.Error(err))
	}

	cfg := &Config{}
	cfg.JWT.AccessSecret = os.Getenv("JWT_ACCESS_SECRET")
	cfg.JWT.RefreshSecret = os.Getenv("JWT_REFRESH_SECRET")
	cfg.Redis.Host = os.Getenv("REDIS_ADDR")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.Server.Port = os.Getenv("PORT")
	return cfg
}
