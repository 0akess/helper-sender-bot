package gitworker

import (
	"context"
	"helper-sender-bot/internal/entity"
	"log/slog"
)

type clientMM interface {
	CreatePost(ctx context.Context, channelID, msg, rootID string) (string, int, error)
}

type repo interface {
	GetAllGitlabConfigs(ctx context.Context) ([]entity.GitlabConfig, error)
	UpdatePostGitMRPushed(ctx context.Context, gitURL string, projectID, mrID int) error
	GetListPostGitMR(ctx context.Context, team, channel string, gitProjectID int) ([]entity.PostGitMR, error)
}

type GitWorker struct {
	log      *slog.Logger
	clientMM clientMM
	repo     repo
}

func NewGitWorker(logger *slog.Logger, clientMM clientMM, repo repo) *GitWorker {
	return &GitWorker{
		log:      logger,
		clientMM: clientMM,
		repo:     repo,
	}
}
