package gitworker

import (
	"context"
	"fmt"
	"time"
)

func (gw *GitWorker) SendDayPing(ctx context.Context) {
	layout := "2006-01-02"

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

			if post.UpdateAT.Format(layout) == time.Now().Format(layout) {
				continue
			}

			msgDayPing := fmt.Sprintf("%s \nДень прошел число сменилось, ничего не изменилось", post.Reviewers)
			_, _, err := gw.clientMM.CreatePost(ctx, post.ChannelID, msgDayPing, post.PostID)
			if err != nil {
				gw.log.Error("Ошибка создания поста в MM", "error", err)
				continue
			}

			err = gw.repo.UpdatePostGitMRPushed(ctx, post.GitURL, post.GitProjectID, post.GitMRID)
			if err != nil {
				gw.log.Warn("Ошибка обновления записи в posts_git_mr", "error", err)
				continue
			}
		}
	}
}
