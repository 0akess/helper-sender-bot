package api

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/controller/api/validator"
	"log/slog"
	"net/http"
	"time"
)

func InitEcho(log *slog.Logger, timeOut time.Duration) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &validator.CustomValidator{Validator: validator.NewValidator()}
	e.HTTPErrorHandler = middleware.HTTPErrorHandler()
	e.Use(
		slogecho.New(log),
		echomw.Recover(),
		echomw.Logger(),
		echomw.RequestID(),
		echomw.ContextTimeoutWithConfig(echomw.ContextTimeoutConfig{Timeout: timeOut}),
	)
	return e
}

func Run(log *slog.Logger, e *echo.Echo, port string) error {
	addr := fmt.Sprintf(":%s", port)
	log.Info("starting server", "addr", addr)
	err := e.Start(addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error run: %w", err)
	}
	return nil
}
