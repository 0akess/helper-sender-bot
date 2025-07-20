package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Logger struct {
	LogLevel string `envconfig:"LOG_LEVEL" required:"true"`
}

func LoadLoggerConfig() (*Logger, error) {
	var cfg Logger

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("env processing: %w", err)
	}
	return &cfg, nil
}
