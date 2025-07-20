package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
	_ "go.uber.org/automaxprocs/maxprocs"
	r "helper-sender-bot/internal/adapters/dbhesebo"
	"helper-sender-bot/internal/adapters/dbhesebo/migrations"
	"helper-sender-bot/internal/applications/config"
	"helper-sender-bot/internal/logger"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	cfg, err := config.LoadLoggerConfig()
	if err != nil {
		panic(fmt.Errorf("load logger config: %w", err))
	}

	appLogger := logger.New(cfg.LogLevel)

	cfgP, err := config.LoadPostgresConfig()
	if err != nil {
		cancel()
		appLogger.Error("Failed load postgres config", "err", err)
		os.Exit(1)
	}

	pool, err := r.NewDB(ctx, cfgP, appLogger)
	if err != nil {
		cancel()
		appLogger.Error("Failed to initialize database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)

	if err = goose.SetDialect("postgres"); err != nil {
		pool.Close()
		cancel()
		appLogger.Error("Failed to set dialect", "err", err)
		os.Exit(1)
	}

	goose.SetBaseFS(migrations.FS)
	parts := strings.Fields(cfgP.Command)
	cmd := parts[0]
	args := parts[1:]

	err = goose.RunContext(ctx, cmd, db, ".", args...)
	if err != nil {
		pool.Close()
		cancel()
		appLogger.Error("Goose command failed", "err", err)
		os.Exit(1)
	}
	appLogger.Info("Goose command executed successfully", "command", cfgP.Command)
}
