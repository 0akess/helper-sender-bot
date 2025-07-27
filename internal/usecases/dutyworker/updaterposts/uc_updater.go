package updaterposts

import (
	"context"
	"helper-sender-bot/internal/entity"
	"log/slog"
	"time"
)

type HTTPClient interface {
	FetchPosts(ctx context.Context, path string) ([]entity.Post, error)
}

type repo interface {
	GetListCfgDuty(ctx context.Context) ([]entity.Chat, error)
	CreatePostDuty(ctx context.Context, channelID, postID string, createdAt time.Time, inProgress bool) error
	DeletePostDuty(ctx context.Context, channelID, postID string) error
}

type PostInfo struct {
	client HTTPClient
	repo   repo
	log    *slog.Logger
}

func NewUpdaterPostInfo(client HTTPClient, repo repo, logger *slog.Logger) *PostInfo {
	return &PostInfo{
		client: client,
		repo:   repo,
		log:    logger,
	}
}

// UpdaterPosts забирает каналы и для каждого выгружает новые посты
func (pi *PostInfo) UpdaterPosts(ctx context.Context, intervalCycle time.Duration) {
	cutoff := time.Now().Add(-intervalCycle).UnixMilli()
	channels, err := pi.repo.GetListCfgDuty(ctx)
	if err != nil {
		pi.log.Error("GetListCfgDuty", "err", err)
		return
	}

	for _, cfg := range channels {
		if cfg.IsNotWorkingHours() {
			return
		}

		pi.fetchAndStore(ctx, cfg.ChannelID, cfg, cutoff)
	}

}
