package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Git struct {
	GitApiToken string `envconfig:"GIT_API_TOKEN" required:"true"`
	GitURL      string `envconfig:"GIT_URL" required:"true"`
}

func LoadGitConfig(prefix string) (*Git, error) {
	var cfg Git

	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("env processing: %w", err)
	}
	return &cfg, nil
}
