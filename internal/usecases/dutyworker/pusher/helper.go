package pusher

import (
	c "context"
	"fmt"
	e "helper-sender-bot/internal/entity"
	"time"
)

// getDutyForChannel вытаскивает имя дежурного из Duty или из header.
func (s *Sender) getDutyForChannel(ctx c.Context, channelID string) (string, error) {
	d, err := s.cacheDuty.GetDutyCache(ctx, channelID)
	if err != nil {
		return "", err
	}
	if d == "" {
		return "@here укажите в хедере канала 'Дежурный: @login'", err
	}
	return d, err
}

// handleThread отвечает за логику какой пуш будем пушить
func (s *Sender) handleThread(ctx c.Context, channelID, duty string, cfg e.Chat, dutyInfo e.PostsInfoDuty) error {
	now := time.Now()
	ageSinceCreate := now.Sub(dutyInfo.CreateAt)
	ageSincePush := now.Sub(dutyInfo.LastPushAt)

	if !dutyInfo.InProgress && ageSinceCreate >= (time.Duration(cfg.DutyTtlInMinute)*time.Minute) {
		return s.firstPush(ctx, channelID, duty, cfg, dutyInfo)
	}

	if dutyInfo.InProgress && ageSincePush >= (time.Duration(cfg.DutyRepeatTtlInMinute)*time.Minute) {
		return s.repeatPush(ctx, channelID, duty, cfg, dutyInfo)
	}
	return nil
}

// firstPush отработает единожды, потому что переведет в InProgress true
func (s *Sender) firstPush(ctx c.Context, channelID, duty string, cfg e.Chat, t e.PostsInfoDuty) error {
	msg := fmt.Sprintf("%s, обрати внимание на обращение :%s:", duty, cfg.EmojiStart)
	err := s.sendAndRecord(ctx, channelID, msg, t)
	if err != nil {
		return err
	}

	return s.repo.MarkPostsDutyAsInProgress(ctx, channelID, t.PostID)
}

// repeatPush повторно пушит дежурного с какой-то периодичностью
func (s *Sender) repeatPush(ctx c.Context, channelID, duty string, cfg e.Chat, t e.PostsInfoDuty) error {
	msg := fmt.Sprintf(
		"%s, :%s: прошло %d мин — не забудь разобрать и поставить :%s:",
		duty, cfg.EmojiStart, cfg.DutyRepeatTtlInMinute, cfg.EmojiDone,
	)
	return s.sendAndRecord(ctx, channelID, msg, t)
}

// sendAndRecord шлёт CreatePost и обновляет last_push_at
func (s *Sender) sendAndRecord(ctx c.Context, channelID, msg string, t e.PostsInfoDuty) error {
	_, stCode, err := s.client.CreatePost(ctx, channelID, msg, t.PostID)
	if err != nil {
		if stCode == 400 {
			return s.repo.DeletePostDuty(ctx, channelID, t.PostID)
		}
		return err
	}

	return s.repo.UpdatePushAtPostDuty(ctx, channelID, t.PostID)
}

// isWorkingHours фильтр позволяющий ограничить в какие часы можно слать уведомления
func (s *Sender) isNotWorkingHours(cfg e.Chat) bool {
	msk := time.FixedZone("MSK", 3*60*60)
	h := time.Now().In(msk).Hour()
	return h < cfg.WorkdayStart || h >= cfg.WorkdayEnd
}
