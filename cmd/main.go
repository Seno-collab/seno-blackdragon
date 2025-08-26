package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	_ "seno-blackdragon/docs"
	"seno-blackdragon/internal/api"
	"seno-blackdragon/internal/config"
	"seno-blackdragon/internal/db"
	"seno-blackdragon/internal/store"
	"seno-blackdragon/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger.Init(logger.LoggerConfig{
		Environment:   "development",
		FilePath:      "logs/app",
		ConsolePretty: true,
		DebugToFile:   true,
		RotateDaily:   true,
	})
	defer logger.Close()
	cfg := config.LoadConfig(logger.Log)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	gin.SetMode(gin.ReleaseMode)

	logger.Log.Info("App started",
		zap.String("Module", "main"),
		zap.Int("version", 1),
	)
	db := db.ConnectDatabase(logger.Log, dsn, cfg.DB.Name)
	redis := store.Config{
		Addr:      cfg.Redis.Host,
		Password:  cfg.Redis.Password,
		Databases: store.DBCache,
	}
	cs, err := store.InitRedis(logger.Log, redis)
	if err != nil {
		logger.Log.Warn("redis_init_partial", zap.Error(err))
	}
	defer cs.Close(logger.Log)
	router := api.InitRouter(db, logger.Log, cs, cfg)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}
	go func() {
		logger.Log.Info("Service running at port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server error", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown", zap.Error(err))
	}
	logger.Log.Info("Server exiting")
}
