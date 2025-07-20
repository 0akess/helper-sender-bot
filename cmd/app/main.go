package main

import (
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"helper-sender-bot/internal/adapters/dbhesebo"
	"helper-sender-bot/internal/adapters/gitlab"
	"helper-sender-bot/internal/adapters/mattermost"
	"helper-sender-bot/internal/applications/config"
	"helper-sender-bot/internal/controller/api"
	chatCon "helper-sender-bot/internal/controller/api/api/c_config_duty"
	gitCon "helper-sender-bot/internal/controller/api/api/c_config_gitlab"
	teamCon "helper-sender-bot/internal/controller/api/api/c_team"
	webhook "helper-sender-bot/internal/controller/api/webhook/wh_gitlab"
	"helper-sender-bot/internal/logger"
	"helper-sender-bot/internal/usecases/auth"
	"helper-sender-bot/internal/usecases/cfgduty"
	"helper-sender-bot/internal/usecases/cfggitlab"
	"helper-sender-bot/internal/usecases/gitworker"
	"helper-sender-bot/internal/usecases/team"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	cfgL, err := config.LoadLoggerConfig()
	if err != nil {
		panic(fmt.Errorf("load logger config: %w", err))
	}

	appLogger := logger.New(cfgL.LogLevel)

	cfgP, err := config.LoadPostgresConfig()
	if err != nil {
		cancel()
		appLogger.Error("Failed load postgres config", "err", err)
		os.Exit(1)
	}

	cfgA, err := config.LoadAppConfig()
	if err != nil {
		cancel()
		appLogger.Error("Failed load app config", "err", err)
		os.Exit(1)
	}

	cfgB, err := config.LoadMattermostBaseConfig()
	if err != nil {
		cancel()
		appLogger.Error("Failed load MM config", "err", err)
		os.Exit(1)
	}

	cfgGitLocal, err := config.LoadGitConfig("LOCAL")
	if err != nil {
		cancel()
		appLogger.Error("Failed load LOCAL Git config", "err", err)
		os.Exit(1)
	}

	gitsCfg := []gitlab.GitConfigs{
		{
			BaseURL: cfgGitLocal.GitURL,
			Token:   cfgGitLocal.GitApiToken,
		},
	}

	gitClient, err := gitlab.New(gitsCfg)
	if err != nil {
		cancel()
		appLogger.Error("Failed initialize gitlab client", "err", err)
		os.Exit(1)
	}

	db, err := dbhesebo.NewDB(ctx, cfgP, appLogger)
	if err != nil {
		cancel()
		appLogger.Error("Failed to initialize database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// различные репозитории
	repo := dbhesebo.NewStorage(db)
	mmClient := mattermost.New(cfgB.MattermostBase, cfgB.Token, appLogger)

	// uc для api
	ucAuth := auth.NewAuth(ctx, repo)
	ucChatConfigs := cfgduty.NewDutyCfgCases(ctx, repo)
	ucTeams := team.NewTeamCases(ctx, repo)
	ucGitConfigs := cfggitlab.NewGitCfgCases(ctx, repo)
	ucGitBot := gitworker.NewSender(appLogger, mmClient, gitClient, repo)

	echo := api.InitEcho(appLogger, cfgA.Timeout)

	chatCon.NewControllerCfgDuty(ctx, ucChatConfigs, ucAuth).RegisterRoutes(echo)
	teamCon.NewControllerTeam(ctx, ucTeams, ucAuth).RegisterRoutes(echo)
	gitCon.NewControllerCfgGit(ctx, ucGitConfigs, ucAuth).RegisterRoutes(echo)
	webhook.NewControllerGitlab(ctx, ucGitBot, cfgA.WebhookToken).RegisterRoutes(echo)

	err = api.Run(appLogger, echo, cfgA.Port)
	if err != nil {
		db.Close()
		cancel()
		appLogger.Error("server failed to start", "err", err)
		os.Exit(1)
	}
}
