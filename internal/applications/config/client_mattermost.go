package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Mattermost struct {
	Token          string `envconfig:"MATTERMOST_TOKEN" required:"true"`
	MattermostBase string `envconfig:"MATTERMOST_BASE" required:"true"`
}

func LoadMattermostBaseConfig() (*Mattermost, error) {
	var cfg Mattermost

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("env processing: %w", err)
	}
	return &cfg, nil
}
