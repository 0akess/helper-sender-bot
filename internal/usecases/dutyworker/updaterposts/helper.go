package updaterposts

import (
	"context"
	"fmt"
	"helper-sender-bot/internal/entity"
	"time"
)

// fetchAndStore выкачивает страницы API начиная с since и обрабатывает каждую
func (pi *PostInfo) fetchAndStore(ctx context.Context, channelID string, cfg entity.Chat, since int64) {
	const perPage = 200
	base := fmt.Sprintf("/api/v4/channels/%s/posts", channelID)

	if since > 0 {
		url := fmt.Sprintf("%s?since=%d&page=0&per_page=%d", base, since, perPage)
		batch, err := pi.client.FetchPosts(ctx, url)
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
		url := fmt.Sprintf("%s?page=%d&per_page=%d", base, page, perPage)
		batch, err := pi.client.FetchPosts(ctx, url)
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
