package gitworker

import (
	"context"
	"fmt"
	"time"
)

func (s *Sender) SendRepeatPush(ctx context.Context) {

	gitCfg, err := s.repo.GetAllGitlabConfigs(ctx)
	if err != nil {
		s.log.Warn("Ошибка получения записей gitlab конфигураций", "error", err)
		return
	}

	for _, cfg := range gitCfg {
		postsGitMr, err := s.repo.GetListPostGitMR(ctx, cfg.Team, cfg.ChannelID, cfg.ProjectID)
		if err != nil {
			s.log.Warn("Ошибка получения записей из posts_git_mr", "error", err)
		}

		for _, post := range postsGitMr {
			ageSinceCreate := time.Since(post.CreateAt)
			sla := time.Duration(post.TTLReview.SLA) * time.Minute
			if ageSinceCreate < sla || post.PushedReview {
				continue
			}

			msgRepeatPush := fmt.Sprintf("%s \nМР не продвинулся, не забывайте посмотреть", post.Reviewers)
			_, _, err := s.clientMM.CreatePost(ctx, post.ChannelID, msgRepeatPush, post.PostID)
			if err != nil {
				s.log.Warn("Ошибка создания поста в MM", "err", err)
				continue
			}

			err = s.repo.UpdatePostGitMRPushed(ctx, post.GitURL, post.GitProjectID, post.GitMRID)
			if err != nil {
				s.log.Warn("Ошибка обновления записи в posts_git_mr", "error", err)
				continue
			}
		}
	}
}
