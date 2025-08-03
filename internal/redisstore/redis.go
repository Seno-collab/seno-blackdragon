package redisstore

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient map[string]*redis.Client

func InitRedis(logger *zap.Logger, hostRedis string, passwordRedis string) RedisClient {
	clients := make(RedisClient)
	configs := []struct {
		Name string
		DB   int
	}{
		{
			Name: "token",
			DB:   0,
		},
	}
	for _, cfg := range configs {
		client := redis.NewClient(&redis.Options{
			Addr:         hostRedis,
			Password:     passwordRedis,
			DB:           cfg.DB,
			PoolSize:     20,
			MinIdleConns: 5,
			PoolTimeout:  30 * time.Second,
		})
		if _, err := client.Ping(context.Background()).Result(); err != nil {
			logger.Panic(fmt.Sprintf("Redis %s connect error: %v", cfg.Name, err))
		}
		clients[cfg.Name] = client
	}
	return clients
}

func CloseAll(logger *zap.Logger, clients RedisClient) {
	for name, c := range clients {
		_ = c.Close()
		logger.Info("Close", zap.String("Redis :", name))
	}
}
