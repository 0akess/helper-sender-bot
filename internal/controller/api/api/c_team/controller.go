package c_chat_config

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/entity"
)

type team interface {
	CreateTeam(ctx context.Context, team entity.Team) error
	GetTeams(ctx context.Context, limit, cursor int, search string) ([]string, int, error)
	UpdateTeam(ctx context.Context, newTeam entity.Team, teamName string, token uuid.UUID) error
	DeleteTeam(ctx context.Context, teamName string, token uuid.UUID) error
}

type auth interface {
	CheckAuth(ctx context.Context, auth entity.AuthMeta) error
}

type Controller struct {
	ctx  context.Context
	team team
	auth auth
}

func NewControllerTeam(ctx context.Context, team team, ucAuth auth) *Controller {
	return &Controller{
		ctx:  ctx,
		team: team,
		auth: ucAuth,
	}
}

func (t *Controller) RegisterRoutes(e *echo.Echo) {
	open := e.Group("/v1")
	open.POST("/team", t.createTeam)
	open.GET("/team", t.getTeam)

	auth := e.Group("/v1")
	auth.Use(middleware.HeaderAuth)
	auth.PUT("/team", t.updateTeam)
	auth.DELETE("/team", t.deleteTeam)
}
