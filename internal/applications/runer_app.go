package applications

import (
	"context"
	"errors"
	_ "github.com/joho/godotenv/autoload"
	_ "go.uber.org/automaxprocs/maxprocs"
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
	"net/http"
	"os/signal"
	"syscall"
)

func RunApp(cfgL *config.Logger) int {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	appLogger := logger.New(cfgL.LogLevel)
	defer appLogger.Info("")

	cfgP, err := config.LoadPostgresConfig()
	if err != nil {
		appLogger.Error("Failed load postgres config", "err", err)
		return 1
	}

	cfgA, err := config.LoadAppConfig()
	if err != nil {
		appLogger.Error("Failed load app config", "err", err)
		return 1
	}

	cfgB, err := config.LoadMattermostBaseConfig()
	if err != nil {
		appLogger.Error("Failed load MM config", "err", err)
		return 1
	}

	cfgGitLocal, err := config.LoadGitConfig("LOCAL")
	if err != nil {
		appLogger.Error("Failed load LOCAL Git config", "err", err)
		return 1
	}

	gitsCfg := []gitlab.GitConfigs{
		{
			BaseURL: cfgGitLocal.GitURL,
			Token:   cfgGitLocal.GitApiToken,
		},
	}

	gitClient, err := gitlab.New(gitsCfg)
	if err != nil {
		appLogger.Error("Failed initialize gitlab client", "err", err)
		return 1
	}

	db, err := dbhesebo.NewDB(ctx, cfgP, appLogger)
	if err != nil {
		appLogger.Error("Failed to initialize database", "err", err)
		return 1
	}

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

	chatCon.NewControllerCfgDuty(ucChatConfigs, ucAuth).RegisterRoutes(echo)
	teamCon.NewControllerTeam(ucTeams, ucAuth).RegisterRoutes(echo)
	gitCon.NewControllerCfgGit(ucGitConfigs, ucAuth).RegisterRoutes(echo)
	webhook.NewControllerGitlab(ctx, ucGitBot, cfgA.WebhookToken).RegisterRoutes(echo)

	go func() {
		if err = api.Run(appLogger, echo, cfgA.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error("server failed to start", "err", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, cfgA.Timeout)
	defer shutdownCancel()

	if err = echo.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("server failed to shutdown", "err", err)
		return 1
	}
	appLogger.Info("server shutdown")

	db.Close()
	return 0
}
