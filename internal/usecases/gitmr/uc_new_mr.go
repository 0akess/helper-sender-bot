package gitmr

import (
	"context"
	"fmt"
	"helper-sender-bot/internal/entity"
)

func (gm *GitMR) SendPushNewMR(ctx context.Context, mr entity.MergeRequestPayload) {
	gitCfg, gitUrl, err := gm.getGitCfg(ctx, mr)
	if err != nil {
		gm.log.Error("Ошибка получения конфигурации", "error", err)
		return
	}

	exist, err := gm.repo.ExistsPostGitMR(ctx, gitUrl, mr.ProjectID, mr.MrID)
	if err != nil {
		gm.log.Error("Ошибка при проверки существования записи", "error", err, "mrID", mr.MrID)
		return
	}
	if !exist && mr.IsDraft {
		return
	}
	if exist && mr.IsDraft {
		err = gm.repo.UpdatePostGitMRIsDraft(ctx, gitUrl, mr.ProjectID, mr.MrID, mr.IsDraft)
		if err != nil {
			gm.log.Error("Ошибка обновления is_draft", "error", err)
			return
		}
		return
	}

	if exist && !mr.IsDraft {
		post, err := gm.repo.GetPostGitMR(ctx, gitUrl, mr.ProjectID, mr.MrID)
		if err != nil {
			gm.log.Error("Ошибка получения posts_git_mr", "error", err)
			return
		}
		if post.IsDraft {
			err = gm.repo.UpdatePostGitMRIsDraft(ctx, gitUrl, mr.ProjectID, mr.MrID, mr.IsDraft)
			if err != nil {
				gm.log.Error("Ошибка обновления is_draft", "error", err)
				return
			}
			msg := fmt.Sprintf("MR вышел из состояния драфта, можно смотреть %s", post.Reviewers)
			_, _, err = gm.clientMM.CreatePost(ctx, gitCfg.ChannelID, msg, post.PostID)
			if err != nil {
				gm.log.Error("CreatePost при MR != Draft", "error", err)
				return
			}
		}
		return
	}

	var mrInfo entity.MergeRequestInfo
	if gitCfg.PushQaAfterReview || len(gitCfg.TTLReview) > 0 {
		mrInfo, err = gm.clientG.GetMRInfo(ctx, mr, gitUrl)
		if err != nil {
			return
		}
	}

	msg, ttlRule, reviewers := buildNewMRMessage(mr, gitCfg, mrInfo)

	id, _, err := gm.clientMM.CreatePost(ctx, gitCfg.ChannelID, msg, "")
	if err != nil {
		gm.log.Error("Ошибка создания поста в MM", "err", err)
		return
	}
	postMR := entity.PostGitMR{
		TeamName:     gitCfg.Team,
		ChannelID:    gitCfg.ChannelID,
		GitURL:       gitUrl,
		GitProjectID: mr.ProjectID,
		GitMRID:      mr.MrID,
		PostID:       id,
		TTLReview:    ttlRule,
		Reviewers:    reviewers,
	}
	err = gm.repo.CreatePostGitMR(ctx, postMR)
	if err != nil {
		gm.log.Error("Ошибка создания записи в posts_git_mr", "err", err)
		return
	}
}
