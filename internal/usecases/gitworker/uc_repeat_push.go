package gitworker

import (
	"context"
	"fmt"
	"time"
)

func (gw *GitWorker) SendRepeatPush(ctx context.Context) {

	gitCfg, err := gw.repo.GetAllGitlabConfigs(ctx)
	if err != nil {
		gw.log.Error("Ошибка получения записей gitlab конфигураций", "error", err)
		return
	}

	for _, cfg := range gitCfg {
		postsGitMr, err := gw.repo.GetListPostGitMR(ctx, cfg.Team, cfg.ChannelID, cfg.ProjectID)
		if err != nil {
			gw.log.Error("Ошибка получения записей из posts_git_mr", "error", err)
		}

		for _, post := range postsGitMr {
			if post.IsDraft {
				continue
			}
			ageSinceCreate := time.Since(post.UpdateAT)
			sla := time.Duration(post.TTLReview.SLA) * time.Minute
			if ageSinceCreate < sla || post.PushedReview {
				continue
			}

			msgRepeatPush := fmt.Sprintf("%s \nМР не продвинулся, не забывайте посмотреть", post.Reviewers)
			_, _, err := gw.clientMM.CreatePost(ctx, post.ChannelID, msgRepeatPush, post.PostID)
			if err != nil {
				gw.log.Error("Ошибка создания поста в MM", "err", err)
				continue
			}

			err = gw.repo.UpdatePostGitMRPushed(ctx, post.GitURL, post.GitProjectID, post.GitMRID)
			if err != nil {
				gw.log.Error("Ошибка обновления записи в posts_git_mr", "error", err)
				continue
			}
		}
	}
}
