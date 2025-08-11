package redisstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// ClientSet wraps named redis clients.
type ClientSet struct {
	clients map[string]*redis.Client
}

// Config defines how to connect and which logical DBs to init.
type Config struct {
	Addr         string        // "host:port"
	Password     string        // optional
	PoolSize     int           // default 20
	MinIdleConns int           // default 5
	DialTimeout  time.Duration // default 5s
	ReadTimeout  time.Duration // default 3s
	WriteTimeout time.Duration // default 3s
	PoolTimeout  time.Duration // default 30s
	// List of named DBs to create clients for.
	Databases []DBConfig
}

// DBConfig describes one named logical DB (0..15 by default).
type DBConfig struct {
	Name string // e.g. "token"
	DB   int    // e.g. 0
}

// Init creates and health-checks all redis clients.
// It returns a ClientSet or an error if any client fails to connect.
func Init(logger *zap.Logger, cfg Config) (*ClientSet, error) {
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 20
	}
	if cfg.MinIdleConns <= 0 {
		cfg.MinIdleConns = 5
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5 * time.Second
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 3 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 3 * time.Second
	}
	if cfg.PoolTimeout == 0 {
		cfg.PoolTimeout = 30 * time.Second
	}
	if len(cfg.Databases) == 0 {
		cfg.Databases = []DBConfig{{Name: "token", DB: 0}}
	}

	cs := &ClientSet{clients: make(map[string]*redis.Client, len(cfg.Databases))}
	var firstErr error

	for _, d := range cfg.Databases {
		rc := redis.NewClient(&redis.Options{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           d.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			PoolTimeout:  cfg.PoolTimeout,
		})

		// Ping with timeout to avoid hanging on startup
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := rc.Ping(ctx).Err(); err != nil {
			cancel()
			logger.Error("redis_connect_failed",
				zap.String("name", d.Name),
				zap.Int("db", d.DB),
				zap.String("addr", cfg.Addr),
				zap.Error(err),
			)
			_ = rc.Close()
			if firstErr == nil {
				firstErr = fmt.Errorf("redis %q (db=%d) connect error: %w", d.Name, d.DB, err)
			}
			continue
		}
		cancel()

		logger.Info("redis_connected",
			zap.String("name", d.Name),
			zap.Int("db", d.DB),
			zap.String("addr", cfg.Addr),
		)
		cs.clients[d.Name] = rc
	}

	if len(cs.clients) == 0 {
		if firstErr == nil {
			firstErr = errors.New("no redis clients initialized")
		}
		return nil, firstErr
	}
	return cs, firstErr
}

// Get returns a named client and whether it exists.
func (c *ClientSet) Get(name string) (*redis.Client, bool) {
	if c == nil {
		return nil, false
	}
	rc, ok := c.clients[name]
	return rc, ok
}

// MustGet returns a named client or panics (useful for wiring in modules that require it).
func (c *ClientSet) MustGet(name string) *redis.Client {
	rc, ok := c.Get(name)
	if !ok {
		panic("redis client not found: " + name)
	}
	return rc
}

// Close closes all clients. Logs per-client result.
func (c *ClientSet) Close(logger *zap.Logger) {
	if c == nil {
		return
	}
	for name, rc := range c.clients {
		if err := rc.Close(); err != nil {
			logger.Warn("redis_close_failed", zap.String("name", name), zap.Error(err))
			continue
		}
		logger.Info("redis_closed", zap.String("name", name))
	}
}
