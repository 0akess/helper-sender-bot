package pusher

import (
	"context"
	"helper-sender-bot/internal/entity"
	"log/slog"
)

type repo interface {
	GetListCfgDuty(ctx context.Context) ([]entity.Chat, error)
	MarkPostsDutyAsInProgress(ctx context.Context, channelID, postID string) error
	UpdatePushAtPostDuty(ctx context.Context, channelID, postID string) error
	DeletePostDuty(ctx context.Context, channelID, postID string) error
	GetListOpenPostDuty(ctx context.Context, channelID string) ([]entity.PostsInfoDuty, error)
}

type client interface {
	CreatePost(ctx context.Context, channelID, msg, rootID string) (string, int, error)
	ChannelHeader(ctx context.Context, id string) (string, error)
}

type cacheDuty interface {
	GetDutyCache(ctx context.Context, channelID string) (string, error)
}

type Sender struct {
	repo      repo
	client    client
	log       *slog.Logger
	cacheDuty cacheDuty
}

func NewPusherDuty(r repo, c client, log *slog.Logger, cDuty cacheDuty) *Sender {
	return &Sender{
		repo:      r,
		client:    c,
		log:       log,
		cacheDuty: cDuty,
	}
}

// PusherBot запускает цикл обработки постов на соблюдения SLA
func (s *Sender) PusherBot(ctx context.Context) {
	channels, err := s.repo.GetListCfgDuty(ctx)
	if err != nil {
		s.log.Error("GetListCfgDuty", "err", err)
		return
	}

	for _, cfg := range channels {
		if cfg.IsNotWorkingHours() {
			return
		}

		threads, err := s.repo.GetListOpenPostDuty(ctx, cfg.ChannelID)
		if err != nil {
			s.log.Warn("list open posts", "channel", cfg.ChannelID, "err", err)
			continue
		}
		if len(threads) == 0 {
			s.log.Info("nothing to push", "channel", cfg.ChannelID)
			continue
		}

		for _, thread := range threads {
			duty, err := s.getDutyForChannel(ctx, cfg.ChannelID)
			if err != nil {
				s.log.Warn("get duty for channel", "channel", cfg.ChannelID, "err", err)
				continue
			}

			err = s.handleThread(ctx, cfg.ChannelID, duty, cfg, thread)
			if err != nil {
				s.log.Error("Failed to handle thread", "channelID", cfg.ChannelID, "err", err)
			}
		}
	}
}
