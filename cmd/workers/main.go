package main

import (
	"fmt"
	"helper-sender-bot/internal/applications"
	"helper-sender-bot/internal/applications/config"
	"os"
)

func main() {
	cfg, err := config.LoadLoggerConfig()
	if err != nil {
		panic(fmt.Errorf("load logger config: %w", err))
	}

	os.Exit(applications.RunWorker(cfg))
}
