package c_config_gitlab

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"helper-sender-bot/internal/entity"
	"net/http"
)

type createGitlabConfig struct {
	GitlabURL         string          `json:"gitlab_url" validate:"required"`
	ProjectName       string          `json:"project_name" validate:"required"`
	ProjectID         int             `json:"project_id" validate:"required"`
	ChannelID         string          `json:"channel_id" validate:"required"`
	Reviewers         []string        `json:"reviewers" validate:"required,min=1,dive,required"`
	ReviewersCount    int             `json:"reviewers_count" validate:"required,min=1"`
	TTLReview         []ttlReviewItem `json:"ttl_review" validate:"required,min=1,dive"`
	QAReviewers       string          `json:"qa_reviewers,omitempty"`
	RequiresQaReview  bool            `json:"requires_qa_review,omitempty"`
	PushQaAfterReview bool            `json:"push_qa_after_review,omitempty"`
}

func (c *Controller) createGitCfg(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return r.NotAuthMassage(err)
	}

	err = c.ucAuth.Auth(c.ctx, auth)
	if err != nil {
		return r.ForbiddenMassage(err)
	}

	var req createGitlabConfig
	if err := e.Bind(&req); err != nil {
		return r.InvalidInputMassage(err)
	}

	if err := e.Validate(req); err != nil {
		return r.InvalidInputMassage(err)
	}

	if (req.RequiresQaReview || req.PushQaAfterReview) && req.QAReviewers == "" {
		return r.InvalidInputMassage(
			fmt.Errorf("qa_reviewers is required for requires_qa_review or push_qa_after_review"),
		)
	}

	ttlReview, err, ok := checkAndBuildTTLReview(req.TTLReview)
	if !ok {
		return r.InvalidInputMassage(err)
	}

	git := entity.GitlabConfig{
		Team:              auth.Team,
		GitlabURL:         req.GitlabURL,
		ProjectName:       req.ProjectName,
		ProjectID:         req.ProjectID,
		ChannelID:         req.ChannelID,
		Reviewers:         req.Reviewers,
		ReviewersCount:    req.ReviewersCount,
		TTLReview:         ttlReview,
		QAReviewers:       req.QAReviewers,
		RequiresQaReview:  req.RequiresQaReview,
		PushQaAfterReview: req.PushQaAfterReview,
	}

	err = c.uc.CreateGitlabConfig(c.ctx, git)
	if err != nil {
		return r.InternalErrorMassage(err)
	}
	return e.JSON(http.StatusCreated, map[string]bool{"success": true})
}
