package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type App struct {
	Timeout      time.Duration `envconfig:"REQ_TIMEOUT" required:"true"`
	Port         string        `envconfig:"API_PORT" required:"true"`
	WebhookToken string        `envconfig:"WEBHOOK_TOKEN" required:"true"`
}

func LoadAppConfig() (*App, error) {
	var cfg App

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("env processing: %w", err)
	}
	return &cfg, nil
}
