package api

import (
	"seno-blackdragon/internal/api/handler"
	"seno-blackdragon/internal/repository"
	"seno-blackdragon/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title           Seno-BlackDragon API
// @version         1.0
// @description     This is a black dragon server.
// @host            localhost:8080
// @BasePath        /api/v1
func InitRouter(db *pgx.Conn, logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := router.Group("/api/v1")
	{
		v1.GET("ping", handler.Ping)
		// auth
		authRepo := repository.NewAuthRepo(db)
		authService := service.NewAuthService(authRepo, logger)
		authHandler := handler.NewAuthHandler(authService)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}
	}
	return router
}
