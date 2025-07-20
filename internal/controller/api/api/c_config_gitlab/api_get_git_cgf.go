package c_config_gitlab

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

type getGitlabConfigRes struct {
	GitlabURL         string          `json:"gitlab_url"`
	ProjectName       string          `json:"project_name"`
	ProjectID         int             `json:"project_id"`
	ChannelID         string          `json:"channel_id"`
	Reviewers         []string        `json:"reviewers"`
	ReviewersCount    int             `json:"reviewers_count"`
	TTLReview         []ttlReviewItem `json:"ttl_review"`
	QAReviewers       string          `json:"qa_reviewers"`
	RequiresQaReview  bool            `json:"requires_qa_review"`
	PushQaAfterReview bool            `json:"push_qa_after_review"`
}

func (c *Controller) getGitCfg(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return r.NotAuthMassage(err)
	}

	err = c.ucAuth.Auth(c.ctx, auth)
	if err != nil {
		return r.ForbiddenMassage(err)
	}

	gitCfg, err := c.uc.GetGitlabConfigs(c.ctx, auth.Team)
	if err != nil {
		return r.InternalErrorMassage(err)
	}

	if len(gitCfg) == 0 {
		return e.JSON(http.StatusOK, []string{})
	}

	res := make([]getGitlabConfigRes, len(gitCfg))
	for i, git := range gitCfg {
		res[i] = getGitlabConfigRes{
			GitlabURL:         git.GitlabURL,
			ProjectName:       git.ProjectName,
			ProjectID:         git.ProjectID,
			ChannelID:         git.ChannelID,
			Reviewers:         git.Reviewers,
			ReviewersCount:    git.ReviewersCount,
			TTLReview:         toGetTTL(git.TTLReview),
			QAReviewers:       git.QAReviewers,
			RequiresQaReview:  git.RequiresQaReview,
			PushQaAfterReview: git.PushQaAfterReview,
		}
	}

	return e.JSON(http.StatusOK, res)
}
