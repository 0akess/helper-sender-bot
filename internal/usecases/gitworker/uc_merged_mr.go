package gitworker

import (
	"context"
	"fmt"
	"helper-sender-bot/internal/entity"
)

func (s *Sender) SendPushMergedMR(ctx context.Context, mr entity.MergeRequestPayload) {
	gitCfg, gitUrl, done := s.getGitCfg(ctx, mr)
	if done {
		return
	}

	postGitMr, err := s.repo.GetPostGitMR(ctx, gitUrl, mr.ProjectID, mr.MRIID)
	if err != nil {
		s.log.Warn("Ошибка получения записи из posts_git_mr", "error", err)
	}

	err = s.repo.DeletePostGitMR(ctx, gitUrl, mr.ProjectID, mr.MRIID)
	if err != nil {
		s.log.Warn("Ошибка удаления записи в posts_git_mr", "error", err)
	}

	msg := "МР влит :large_green_square:"
	if gitCfg.PushQaAfterReview {
		msg = fmt.Sprintf("%s\n\nПризываю к проверке %s", msg, gitCfg.QAReviewers)
	}

	_, _, err = s.clientMM.CreatePost(ctx, gitCfg.ChannelID, msg, postGitMr.PostID)
	if err != nil {
		s.log.Warn("Ошибка создания поста в MM", "err", err)
		return
	}
}
