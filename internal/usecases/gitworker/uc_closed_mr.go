package gitworker

import (
	"context"
	"helper-sender-bot/internal/entity"
)

const msgClosedMR = "МР закрыт :large_red_square:"

func (s *Sender) SendPushClosedMR(ctx context.Context, mr entity.MergeRequestPayload) {
	gitCfg, gitUrl, done := s.getGitCfg(ctx, mr)
	if done {
		return
	}

	postGitMr, err := s.repo.GetPostGitMR(ctx, gitUrl, mr.ProjectID, mr.MRIID)
	if err != nil {
		s.log.Warn("Ошибка получения записи из posts_git_mr", "error", err)
	}

	_, _, err = s.clientMM.CreatePost(ctx, gitCfg.ChannelID, msgClosedMR, postGitMr.PostID)
	if err != nil {
		s.log.Warn("Ошибка создания поста в MM", "err", err)
		return
	}

	err = s.repo.DeletePostGitMR(ctx, gitUrl, mr.ProjectID, mr.MRIID)
	if err != nil {
		s.log.Warn("Ошибка удаления записи в posts_git_mr", "error", err)
	}
}
