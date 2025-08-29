package updaterposts

import (
	"context"
	"helper-sender-bot/internal/entity"
	"time"
)

// fetchAndStore выкачивает страницы API начиная с since и обрабатывает каждую
func (pi *PostInfo) fetchAndStore(ctx context.Context, channelID string, cfg entity.Chat, since int64) {
	const perPage int = 200

	if since > 0 {
		batch, err := pi.client.FetchPostsWithSince(ctx, channelID, int(since), perPage)
		if err != nil {
			pi.log.Error("fetch posts", "channel", channelID, "err", err)
			return
		}

		var filtered []entity.Post
		for _, p := range batch {
			if p.CreateAt > since {
				filtered = append(filtered, p)
			}
		}
		if len(filtered) > 0 {
			pi.processBatch(ctx, channelID, cfg, filtered)
		}
		return
	}

	for page := 0; ; page++ {
		batch, err := pi.client.FetchPostsByPage(ctx, channelID, page, perPage)
		if err != nil {
			pi.log.Error("fetch posts", "channel", channelID, "err", err)
			return
		}
		if len(batch) == 0 {
			break
		}

		var filtered []entity.Post
		for _, p := range batch {
			if p.CreateAt > since {
				filtered = append(filtered, p)
			}
		}
		if len(filtered) == 0 {
			break
		}

		pi.processBatch(ctx, channelID, cfg, filtered)

		if len(batch) < perPage {
			break
		}
	}
}

// processBatch сохраняет в БД только топ-левел посты, пропуская системные, внутри-тредовые и закрытые.
func (pi *PostInfo) processBatch(ctx context.Context, channelID string, cfg entity.Chat, batch []entity.Post) {
	for _, p := range batch {
		if p.RootId != "" || p.Type != "" {
			continue
		}
		if p.IsExistReaction(cfg.EmojiDone) {
			err := pi.repo.DeletePostDuty(ctx, channelID, p.ID)
			if err != nil {
				pi.log.Error("delete post", "channel", channelID, "root_id", p.RootId, "err", err)
			}
			continue
		}
		inProgress := p.IsExistReaction(cfg.EmojiStart)

		createdAt := time.UnixMilli(p.CreateAt)

		err := pi.repo.CreatePostDuty(ctx, channelID, p.ID, createdAt, inProgress)
		if err != nil {
			pi.log.Error("CreatePostDuty", "channel", channelID, "post", p.ID, "err", err)
		}
	}
}
