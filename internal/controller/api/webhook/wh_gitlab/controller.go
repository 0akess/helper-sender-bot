package wh_gitlab

import (
	"context"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/entity"
	"log/slog"
)

type pushMR interface {
	SendPushNewMR(ctx context.Context, mr entity.MergeRequestPayload)
	SendPushMergedMR(ctx context.Context, mr entity.MergeRequestPayload)
	SendPushClosedMR(ctx context.Context, mr entity.MergeRequestPayload)
}

type GitlabController struct {
	ctx    context.Context
	pushMR pushMR
	log    *slog.Logger
	token  string
}

func NewControllerGitlab(ctx context.Context, pushMR pushMR, log *slog.Logger, token string) *GitlabController {
	c := &GitlabController{
		ctx:    ctx,
		pushMR: pushMR,
		log:    log,
		token:  token,
	}
	return c
}

func (gc *GitlabController) RegisterRoutes(e *echo.Echo) {
	e.POST("/gitlab/webhook/mr_info", gc.handleGitlabWebhook)
}
