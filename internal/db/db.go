package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func ConnectDatabase(logger *zap.Logger, dsn string, dbName string) *pgx.Conn {
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		logger.Error("Unable to parse DB config", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbConn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		logger.Error("Unable to connect to database", zap.Error(err))
	}
	if err := dbConn.Ping(ctx); err != nil {
		logger.Error("Database ping failed", zap.Error(err))
	}
	logger.Info("Connected to PostgreSQL database", zap.String("db", dbName))
	return dbConn
}
