package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"seno-blackdragon/internal/db"
	"seno-blackdragon/internal/redisstore"
	"seno-blackdragon/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title           Seno-BlackDragon API
// @version         1.0
// @description     This is a black dragon server.
// @host            localhost:8080
// @BasePath        /api/v1
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
	db.ConnectDatabase(logger.Log, dsn, dbName)
	redisClients := redisstore.InitRedis(logger.Log, redisAddr, redisPassword)
	defer redisstore.CloseAll(logger.Log, redisClients)
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// PingExample godoc
	// @Summary      Ping example
	// @Description  Do ping
	// @Tags         ping
	// @Success      200  {object}  map[string]string
	// @Router       /ping [get]
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
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
