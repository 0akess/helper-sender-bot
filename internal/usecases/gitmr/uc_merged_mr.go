package gitmr

import (
	"context"
	"fmt"
	"helper-sender-bot/internal/entity"
)

func (gm *GitMR) SendPushMergedMR(ctx context.Context, mr entity.MergeRequestPayload) {
	gitCfg, gitUrl, err := gm.getGitCfg(ctx, mr)
	if err != nil {
		gm.log.Error("Ошибка получения конфигурации", "error", err)
		return
	}

	postGitMr, err := gm.repo.GetPostGitMR(ctx, gitUrl, mr.ProjectID, mr.MrID)
	if err != nil {
		gm.log.Error("Ошибка получения записи из posts_git_mr", "error", err)
	}

	err = gm.repo.DeletePostGitMR(ctx, gitUrl, mr.ProjectID, mr.MrID)
	if err != nil {
		gm.log.Error("Ошибка удаления записи в posts_git_mr", "error", err)
	}

	msg := "МР влит :large_green_square:"
	if gitCfg.PushQaAfterReview {
		msg = fmt.Sprintf("%s\n\nПризываю к проверке %s", msg, gitCfg.QAReviewers)
	}

	_, _, err = gm.clientMM.CreatePost(ctx, gitCfg.ChannelID, msg, postGitMr.PostID)
	if err != nil {
		gm.log.Error("Ошибка создания поста в MM", "error", err)
		return
	}
}
