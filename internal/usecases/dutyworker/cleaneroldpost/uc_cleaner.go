package cleaneroldpost

import (
	"context"
	"helper-sender-bot/internal/entity"
	"log/slog"
	"time"
)

type repo interface {
	DeletePostDutyOlderThan(ctx context.Context, channelID string, cutoff time.Time) error
	GetListCfgDuty(ctx context.Context) ([]entity.Chat, error)
}

type Clean struct {
	repo repo
	log  *slog.Logger
}

func NewCleanOldPost(repo repo, logger *slog.Logger) *Clean {
	return &Clean{
		repo: repo,
		log:  logger,
	}
}

// CleanerOldPost удаляет записи старше intervalCycle по create_at
func (cl *Clean) CleanerOldPost(ctx context.Context, intervalCycle time.Duration) {
	cutoff := time.Now().Add(-intervalCycle)
	channels, err := cl.repo.GetListCfgDuty(ctx)
	if err != nil {
		cl.log.Error("Failed to get list of duty channels", "error", err)
		return
	}

	for _, channelID := range channels {
		if err := cl.repo.DeletePostDutyOlderThan(ctx, channelID.ChannelID, cutoff); err != nil {
			cl.log.Error("DeletePostDutyOlderThan", "channel", channelID, "err", err)
		}
	}

}
