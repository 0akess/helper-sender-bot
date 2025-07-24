package c_config_gitlab

import (
	"context"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/entity"
)

type gitlabCfg interface {
	CreateGitlabConfig(ctx context.Context, cfg entity.GitlabConfig) error
	DeleteGitlabConfig(ctx context.Context, gitProjectID int, gitURL, team string) error
	UpdateGitlabConfig(ctx context.Context, cfg entity.GitlabConfig, gitProjectID int, gitURL string) error
	GetGitlabConfigs(ctx context.Context, team string) ([]entity.GitlabConfig, error)
}

type auth interface {
	CheckAuth(ctx context.Context, auth entity.AuthMeta) error
}

type Controller struct {
	gitlabCfg gitlabCfg
	auth      auth
}

func NewControllerCfgGit(gitlabCfg gitlabCfg, auth auth) *Controller {
	return &Controller{
		gitlabCfg: gitlabCfg,
		auth:      auth,
	}
}

func (c *Controller) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/v1")
	g.Use(middleware.HeaderAuth)
	g.POST("/config_gitlab", c.createGitCfg)
	g.GET("/config_gitlab", c.getGitCfg)
	g.DELETE("/config_gitlab", c.deleteGitCfg)
	g.PUT("/config_gitlab", c.putGitCfg)
}
