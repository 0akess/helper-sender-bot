package gitworker

import (
	"context"
	"fmt"
	"time"
)

func (s *Sender) SendDayPing(ctx context.Context) {
	layout := "2006-01-02"

	cacheGit, err := s.repo.GetAllGitlabConfigs(ctx)
	if err != nil {
		s.log.Warn("Ошибка получения записей gitlab конфигураций", "error", err)
		return
	}

	for _, cfg := range cacheGit {
		postsGitMr, err := s.repo.GetListPostGitMR(ctx, cfg.Team, cfg.ChannelID)
		if err != nil {
			s.log.Warn("Ошибка получения записей из posts_git_mr", "error", err)
		}

		for _, post := range postsGitMr {

			if post.CreateAt.Format(layout) == time.Now().Format(layout) {
				continue
			}

			msgDayPing := fmt.Sprintf("%s \nДень прошел число сменилось, ничего не изменилось", post.Reviewers)
			_, _, err := s.clientMM.CreatePost(ctx, post.ChannelID, msgDayPing, post.PostID)
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
