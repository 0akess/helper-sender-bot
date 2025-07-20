package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Postgres struct {
	Command  string `envconfig:"MIGRATION_COMMAND" required:"true"`
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"true"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	DBName   string `envconfig:"POSTGRES_DBNAME" required:"true"`
	SSLMode  string `envconfig:"POSTGRES_SSLMODE" required:"true"`
	Conns    int32  `envconfig:"POSTGRES_CONNS" required:"true"`
}

func LoadPostgresConfig() (*Postgres, error) {
	var cfg Postgres

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("env processing: %w", err)
	}
	return &cfg, nil
}
