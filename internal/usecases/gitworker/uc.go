package gitworker

import (
	"context"
	"helper-sender-bot/internal/entity"
	"log/slog"
)

type clientMM interface {
	CreatePost(ctx context.Context, channelID, msg, rootID string) (string, int, error)
}

type clientGit interface {
	GetMRInfo(ctx context.Context, p entity.MergeRequestPayload, baseURL string) (entity.MergeRequestInfo, error)
}

type repo interface {
	GetAllGitlabConfigs(ctx context.Context) ([]entity.GitlabConfig, error)
	GetGitlabConfig(ctx context.Context, projectID int, gitUrl string) (entity.GitlabConfig, error)
	CreatePostGitMR(ctx context.Context, p entity.PostGitMR) error
	ExistsPostGitMR(ctx context.Context, gitURL string, projectID, mrID int) (bool, error)
	DeletePostGitMR(ctx context.Context, gitURL string, projectID, mrID int) error
	GetPostGitMR(ctx context.Context, gitURL string, projectID, mrID int) (entity.PostGitMR, error)
	UpdatePostGitMRPushed(ctx context.Context, gitURL string, projectID, mrID int) error
	GetListPostGitMR(ctx context.Context, team, channel string) ([]entity.PostGitMR, error)
}

type Sender struct {
	log      *slog.Logger
	clientMM clientMM
	clientG  clientGit
	repo     repo
}

func NewSender(logger *slog.Logger, clientMM clientMM, clientG clientGit, repo repo) *Sender {
	return &Sender{
		log:      logger,
		clientMM: clientMM,
		clientG:  clientG,
		repo:     repo,
	}
}
