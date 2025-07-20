package c_config_duty

import (
	"context"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/entity"
)

type usecases interface {
	GetListDutyCfgByTeam(ctx context.Context, team string) ([]entity.Chat, error)
	CreateDutyCfg(ctx context.Context, chat entity.Chat, team string) error
	UpdateDutyCfg(ctx context.Context, channel, team string, upd entity.Chat) error
	DeleteDutyCfg(ctx context.Context, channel, team string) error
}

type ucAuth interface {
	Auth(ctx context.Context, auth entity.AuthMeta) error
}

type CfgDutyController struct {
	ctx    context.Context
	uc     usecases
	ucAuth ucAuth
}

func NewControllerCfgDuty(ctx context.Context, usecases usecases, ucAuth ucAuth) *CfgDutyController {
	c := &CfgDutyController{
		ctx:    ctx,
		uc:     usecases,
		ucAuth: ucAuth,
	}
	return c
}

func (c *CfgDutyController) RegisterRoutes(e *echo.Echo) {

	g := e.Group("/v1")
	g.Use(middleware.HeaderAuth)
	g.POST("/config_duty", c.createChat)
	g.GET("/config_duty", c.getChats)
	g.PUT("/config_duty", c.updateChat)
	g.DELETE("/config_duty", c.deleteChat)
}
