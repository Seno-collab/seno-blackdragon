package main

import (
	"seno-blackdragon/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	logger.Init(logger.LoggerConfig{
		Environment: "production",
		FilePath:    "logs/app",
	})
	defer logger.Close()
	logger.Log.Info("App started",
		zap.String("Module", "main"),
		zap.Int("version", 1),
	)
	logger.Sugar.Infof("Hello %s", "Zap")
}
