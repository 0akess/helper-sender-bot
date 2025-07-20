package wh_gitlab

import (
	"context"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/entity"
)

type usecases interface {
	SendPushNewMR(ctx context.Context, mr entity.MergeRequestPayload)
	SendPushMergedMR(ctx context.Context, mr entity.MergeRequestPayload)
	SendPushClosedMR(ctx context.Context, mr entity.MergeRequestPayload)
}

type GitlabController struct {
	ctx   context.Context
	uc    usecases
	token string
}

func NewControllerGitlab(ctx context.Context, uc usecases, token string) *GitlabController {
	c := &GitlabController{
		ctx:   ctx,
		uc:    uc,
		token: token,
	}
	return c
}

func (gc *GitlabController) RegisterRoutes(e *echo.Echo) {
	e.POST("/gitlab/webhook/mr_info", gc.handleGitlabWebhook)
}
