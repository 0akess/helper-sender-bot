package gitworker

import (
	"context"
	"helper-sender-bot/internal/entity"
)

func (s *Sender) SendPushNewMR(ctx context.Context, mr entity.MergeRequestPayload) {
	gitCfg, gitUrl, done := s.getGitCfg(ctx, mr)
	if done {
		return
	}
	exist, err := s.repo.ExistsPostGitMR(ctx, gitUrl, mr.ProjectID, mr.MRIID)
	if err != nil {
		s.log.Warn("Error checking if repo exists", "error", err)
		return
	}
	if exist {
		return
	}

	var mrInfo entity.MergeRequestInfo
	if gitCfg.PushQaAfterReview || len(gitCfg.TTLReview) > 0 {
		mrInfo, err = s.clientG.GetMRInfo(ctx, mr, gitUrl)
		if err != nil {
			return
		}
	}

	msg, ttlRule, reviewers := buildNewMRMessage(mr, gitCfg, mrInfo)

	id, _, err := s.clientMM.CreatePost(ctx, gitCfg.ChannelID, msg, "")
	if err != nil {
		s.log.Warn("Ошибка создания поста в MM", "err", err)
		return
	}
	postMR := entity.PostGitMR{
		TeamName:     gitCfg.Team,
		ChannelID:    gitCfg.ChannelID,
		GitURL:       gitUrl,
		GitProjectID: mr.ProjectID,
		GitMRID:      mr.MRIID,
		PostID:       id,
		TTLReview:    ttlRule,
		Reviewers:    reviewers,
	}
	err = s.repo.CreatePostGitMR(ctx, postMR)
	if err != nil {
		s.log.Warn("Ошибка создания записи в posts_git_mr", "err", err)
		return
	}
}
