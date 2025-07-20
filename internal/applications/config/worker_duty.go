package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Duty struct {
	// Workers какой период для повторного запуска
	CleanOldPost    time.Duration `envconfig:"WAIT_FOR_RUN_BOT_CLEAN_OLD_POST" required:"true"`
	UpdaterPostInfo time.Duration `envconfig:"WAIT_FOR_RUN_BOT_UPDATER_POST_INFO" required:"true"`
	Pusher          time.Duration `envconfig:"WAIT_FOR_RUN_BOT_PUSHER" required:"true"`

	// Жизненный цикл поста-обращения в чате
	PostLifecycle time.Duration `envconfig:"DUTY_POST_LIFECYCLE_TTL" required:"true"`

	// Cache какой период не обновлять кеши
	CacheDuty time.Duration `envconfig:"CACHE_DUTY_SLA" required:"true"`

	// Запускать или нет worker
	StartDuty bool `envconfig:"START_DUTY_WORKER" default:"true"`
}

func LoadDutyWorkerCfg() (*Duty, error) {
	var cfg Duty

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("env processing: %w", err)
	}
	return &cfg, nil
}
