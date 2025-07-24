package applications

import (
	"context"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
	_ "go.uber.org/automaxprocs/maxprocs"
	r "helper-sender-bot/internal/adapters/dbhesebo"
	"helper-sender-bot/internal/adapters/dbhesebo/migrations"
	"helper-sender-bot/internal/applications/config"
	"helper-sender-bot/internal/logger"
	"os/signal"
	"strings"
	"syscall"
)

func RunMigrate(cfg *config.Logger) int {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	appLogger := logger.New(cfg.LogLevel)

	cfgP, err := config.LoadPostgresConfig()
	if err != nil {
		appLogger.Error("Failed load postgres config", "err", err)
		return 1
	}

	pool, err := r.NewDB(ctx, cfgP, appLogger)
	if err != nil {
		appLogger.Error("Failed to initialize database", "err", err)
		return 1
	}

	db := stdlib.OpenDBFromPool(pool)

	if err = goose.SetDialect("postgres"); err != nil {
		appLogger.Error("Failed to set dialect", "err", err)
		return 1
	}

	goose.SetBaseFS(migrations.FS)
	parts := strings.Fields(cfgP.Command)
	cmd := parts[0]
	args := parts[1:]

	err = goose.RunContext(ctx, cmd, db, ".", args...)
	if err != nil {
		appLogger.Error("Goose command failed", "err", err)
		return 1
	}
	appLogger.Info("Goose command executed successfully", "command", cfgP.Command)
	pool.Close()
	return 0
}
