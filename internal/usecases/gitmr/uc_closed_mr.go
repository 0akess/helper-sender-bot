package gitmr

import (
	"context"
	"helper-sender-bot/internal/entity"
)

func (gm *GitMR) SendPushClosedMR(ctx context.Context, mr entity.MergeRequestPayload) {
	gitCfg, gitUrl, err := gm.getGitCfg(ctx, mr)
	if err != nil {
		gm.log.Error("Ошибка получения конфигурации", "error", err)
		return
	}

	postGitMr, err := gm.repo.GetPostGitMR(ctx, gitUrl, mr.ProjectID, mr.MrID)
	if err != nil {
		gm.log.Error("Ошибка получения записи из posts_git_mr", "error", err)
	}

	_, _, err = gm.clientMM.CreatePost(ctx, gitCfg.ChannelID, "МР закрыт :large_red_square:", postGitMr.PostID)
	if err != nil {
		gm.log.Error("Ошибка создания поста в MM", "err", err)
		return
	}

	err = gm.repo.DeletePostGitMR(ctx, gitUrl, mr.ProjectID, mr.MrID)
	if err != nil {
		gm.log.Error("Ошибка удаления записи в posts_git_mr", "error", err)
	}
}
