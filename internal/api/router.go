package api

import (
	"net/http"
	"seno-blackdragon/internal/api/handler"
	"seno-blackdragon/internal/config"
	"seno-blackdragon/internal/repository"
	"seno-blackdragon/internal/service"
	"seno-blackdragon/internal/store"
	"seno-blackdragon/internal/version"
	"seno-blackdragon/pkg/middleware"
	"seno-blackdragon/pkg/pass"
	"time"

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
func InitRouter(db *pgx.Conn, logger *zap.Logger, redis *store.ClientSet, cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.TraceAndLogFullMiddleware(logger, nil))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/docs/index.html")
	})
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// === Health endpoints
	router.GET("/health", func(c *gin.Context) { // liveness
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"uptime": "alive",
		})
	})
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":   version.Version,
			"commit":    version.Commit,
			"buildTime": version.BuildTime,
		})
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("ping", handler.Ping)
		// auth
		jwtCfg := service.JWTConfig{
			AccessSecret:  []byte(cfg.JwtAccessSecret),
			RefreshSecret: []byte(cfg.JwtRefreshSecret),
			AccessTTL:     15 * time.Minute,
			RefreshTTL:    30 * 24 * time.Hour,
			Issuer:        "seno-blackdragon",
		}
		hasher := pass.NewBcryptHasher(pass.BcryptOptions{Cost: 12})

		authRepo := repository.NewUserRepo(db)
		authService := service.NewAuthService(authRepo, hasher, jwtCfg, logger)
		authHandler := handler.NewAuthHandler(authService)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}
	}
	return router
}
