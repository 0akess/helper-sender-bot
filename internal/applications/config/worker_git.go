package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type GitRunCfg struct {
	Pusher   time.Duration `envconfig:"WAIT_FOR_RUN_GIT_PUSHER" required:"true"`
	StartGit bool          `envconfig:"START_GIT_WORKER" default:"true"`
}

func LoadGitWorkerCfg() (*GitRunCfg, error) {
	var cfg GitRunCfg

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("env processing: %w", err)
	}
	return &cfg, nil
}
