package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var DB *pgxpool.Pool

func ConnectDatabase(logger *zap.Logger, dsn string, dbName string) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("Unable to parse DB config", zap.Error(err))
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = 30 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Error("Unable to connect to database", zap.Error(err))
	}
	if err := dbpool.Ping(ctx); err != nil {
		logger.Error("Database ping failed", zap.Error(err))
	}
	logger.Info("Connected to PostgreSQL database", zap.String("db", dbName))
	DB = dbpool
}
