package dbhesebo

import (
	"context"
	"fmt"
	slogadapter "github.com/induzo/gocom/database/pgx-slog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"helper-sender-bot/internal/applications/config"
	"log/slog"
	"net"
	"strconv"
	"time"
)

func NewDB(ctx context.Context, cfg *config.Postgres, logger *slog.Logger) (*pgxpool.Pool, error) {
	hostPort := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User, cfg.Password, hostPort, cfg.DBName, cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	poolConfig.MaxConns = cfg.Conns

	poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   slogadapter.NewLogger(logger),
		LogLevel: tracelog.LogLevelInfo,
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			stats := pool.Stat()
			logger.Info("pgxpool stats",
				slog.Int("acquired", int(stats.AcquiredConns())),
				slog.Int("idle", int(stats.IdleConns())),
				slog.Int("total", int(stats.TotalConns())),
				slog.Int("max", int(stats.MaxConns())),
			)
		}
	}()

	return pool, nil
}
