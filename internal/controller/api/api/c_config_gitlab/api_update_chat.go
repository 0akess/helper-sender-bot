package c_config_gitlab

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/controller/api/api/responses"
	"helper-sender-bot/internal/entity"
	"net/http"
)

type putGitlabConfig struct {
	Reviewers         []string        `json:"reviewers" validate:"required"`
	ReviewersCount    int             `json:"reviewers_count" validate:"required"`
	TTLReview         []ttlReviewItem `json:"ttl_review,omitempty"`
	QAReviewers       string          `json:"qa_reviewers,omitempty"`
	RequiresQaReview  bool            `json:"requires_qa_review,omitempty"`
	PushQaAfterReview bool            `json:"push_qa_after_review,omitempty"`
}

func (c *Controller) putGitCfg(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return responses.NotAuthMessage(err)
	}

	err = c.auth.CheckAuth(e.Request().Context(), auth)
	if err != nil {
		return responses.ForbiddenMessage(err)
	}

	gitURL, projectID, err := gitUrlAndIdQuery(e)
	if err != nil {
		return responses.InvalidInputMessage(err)
	}

	var req putGitlabConfig
	if err := e.Bind(&req); err != nil {
		return responses.InvalidInputMessage(err)
	}
	if err := e.Validate(&req); err != nil {
		return responses.InvalidInputMessage(err)
	}

	var ttlReview []entity.TTLReviewItem
	if len(req.TTLReview) > 0 {
		ttlReview, err = checkAndBuildTTLReview(req.TTLReview)
		if err != nil {
			return responses.InvalidInputMessage(err)
		}
	}

	git := entity.GitlabConfig{
		Team:              auth.Team,
		Reviewers:         req.Reviewers,
		ReviewersCount:    req.ReviewersCount,
		TTLReview:         ttlReview,
		QAReviewers:       req.QAReviewers,
		RequiresQaReview:  req.RequiresQaReview,
		PushQaAfterReview: req.PushQaAfterReview,
	}

	err = c.gitlabCfg.UpdateGitlabConfig(e.Request().Context(), git, projectID, gitURL)
	if err != nil {
		return responses.InternalErrorMessage(err)
	}
	return e.NoContent(http.StatusOK)
}
