package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	_ "seno-blackdragon/docs"
	"seno-blackdragon/internal/api"
	"seno-blackdragon/internal/db"
	"seno-blackdragon/internal/redisstore"
	"seno-blackdragon/pkg/logger"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger.Init(logger.LoggerConfig{
		Environment: "production",
		FilePath:    "logs/app",
	})
	defer logger.Close()
	if err := godotenv.Load(); err != nil {
		logger.Log.Error("Count not load .env file", zap.Error(err))
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	portServer := os.Getenv("PORT")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	logger.Log.Info("App started",
		zap.String("Module", "main"),
		zap.Int("version", 1),
	)
	db := db.ConnectDatabase(logger.Log, dsn, dbName)
	redisClients := redisstore.InitRedis(logger.Log, redisAddr, redisPassword)
	defer redisstore.CloseAll(logger.Log, redisClients)
	router := api.InitRouter(db, logger.Log)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", portServer),
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
