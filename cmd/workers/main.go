package main

import (
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"helper-sender-bot/internal/adapters/cacheduty"
	"helper-sender-bot/internal/adapters/dbhesebo"
	"helper-sender-bot/internal/adapters/gitlab"
	"helper-sender-bot/internal/adapters/mattermost"
	"helper-sender-bot/internal/applications/config"
	cleanerRunner "helper-sender-bot/internal/controller/workers/duty/cleaner_old_post"
	pusherRunner "helper-sender-bot/internal/controller/workers/duty/pusher"
	updaterRunner "helper-sender-bot/internal/controller/workers/duty/updater_posts"
	dayPingerRunner "helper-sender-bot/internal/controller/workers/git/daypinger"
	gitPingSlaRunner "helper-sender-bot/internal/controller/workers/git/pingonsla"
	"helper-sender-bot/internal/logger"
	"helper-sender-bot/internal/usecases/dutyworker/cleaneroldpost"
	"helper-sender-bot/internal/usecases/dutyworker/pusher"
	"helper-sender-bot/internal/usecases/dutyworker/updaterposts"
	"helper-sender-bot/internal/usecases/gitworker"
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

	cfgB, err := config.LoadMattermostBaseConfig()
	if err != nil {
		cancel()
		appLogger.Error("Failed load MM config", "err", err)
		os.Exit(1)
	}

	dutyWorker, err := config.LoadDutyWorkerCfg()
	if err != nil {
		cancel()
		appLogger.Error("Failed load Bot config", "err", err)
		os.Exit(1)
	}

	cfgGitLocal, err := config.LoadGitConfig("LOCAL")
	if err != nil {
		cancel()
		appLogger.Error("Failed load LOCAL Git config", "err", err)
		os.Exit(1)
	}

	gitWorker, err := config.LoadGitWorkerCfg()
	if err != nil {
		cancel()
		appLogger.Error("Failed load LOCAL Git config", "err", err)
		os.Exit(1)
	}

	db, err := dbhesebo.NewDB(ctx, cfgP, appLogger)
	if err != nil {
		cancel()
		appLogger.Error("Failed to initialize database", "err", err)
		os.Exit(1)
	}
	defer db.Close()
	repo := dbhesebo.NewStorage(db)

	gitsCfg := []gitlab.GitConfigs{
		{
			BaseURL: cfgGitLocal.GitURL,
			Token:   cfgGitLocal.GitApiToken,
		},
	}

	mmClient := mattermost.New(cfgB.MattermostBase, cfgB.Token, appLogger)

	// запуск worker для обработки сценариев gitlab
	if gitWorker.StartGit {
		gitClient, err := gitlab.New(gitsCfg)
		if err != nil {
			db.Close()
			cancel()
			appLogger.Error("Failed initialize gitlab client", "err", err)
			os.Exit(1)
		}

		ucGitBot := gitworker.NewSender(appLogger, mmClient, gitClient, repo)

		gitPingSlaRunner.NewRepeatPush(ucGitBot).RunGoSendRepeatPush(ctx, gitWorker.Pusher)
		dayPingerRunner.NewDayPinger(ucGitBot).RunGoSendDayPinger(ctx)
	}

	// запуск worker для обработки сценариев дежурств
	if dutyWorker.StartDuty {
		cacheDuty := cacheduty.NewCache(dutyWorker.CacheDuty, mmClient, appLogger)

		ucPostInfo := updaterposts.NewUpdaterPostInfo(mmClient, repo, appLogger)
		ucCleanOld := cleaneroldpost.NewCleanOldPost(repo, appLogger)
		ucPusher := pusher.NewPusherDuty(repo, mmClient, appLogger, cacheDuty)

		cleanerRunner.NewCleaner(ucCleanOld).
			RunGoCleanerOldPost(ctx, dutyWorker.PostLifecycle, dutyWorker.CleanOldPost)
		updaterRunner.NewUpdaterPosts(ucPostInfo).
			RunGoUpdaterPosts(ctx, dutyWorker.PostLifecycle, dutyWorker.UpdaterPostInfo)
		pusherRunner.NewPusher(ucPusher).
			RunGoPusherBot(ctx, dutyWorker.Pusher)
	}
	<-ctx.Done()

	appLogger.Info("Worker started")
}
