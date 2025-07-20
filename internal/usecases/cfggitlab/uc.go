package cfggitlab

import (
	c "context"
	e "helper-sender-bot/internal/entity"
)

type repo interface {
	CreateGitlabConfig(ctx c.Context, cfg e.GitlabConfig) error
	DeleteGitlabConfigByProjectID(ctx c.Context, gitProjectID int, gitURL, team string) error
	UpdateGitlabConfig(ctx c.Context, cfg e.GitlabConfig, gitProjectID int, gitURL string) error
	GetGitlabConfigsByTeam(ctx c.Context, team string) ([]e.GitlabConfig, error)
}

type GitCfgCases struct {
	ctx  c.Context
	repo repo
}

func NewGitCfgCases(ctx c.Context, repo repo) *GitCfgCases {
	return &GitCfgCases{
		ctx:  ctx,
		repo: repo,
	}
}

func (g *GitCfgCases) CreateGitlabConfig(ctx c.Context, cfg e.GitlabConfig) error {
	return g.repo.CreateGitlabConfig(ctx, cfg)
}

func (g *GitCfgCases) DeleteGitlabConfig(ctx c.Context, projectID int, gitURL, team string) error {
	return g.repo.DeleteGitlabConfigByProjectID(ctx, projectID, gitURL, team)
}

func (g *GitCfgCases) UpdateGitlabConfig(ctx c.Context, cfg e.GitlabConfig, gitProjectID int, gitURL string) error {
	return g.repo.UpdateGitlabConfig(ctx, cfg, gitProjectID, gitURL)
}

func (g *GitCfgCases) GetGitlabConfigs(ctx c.Context, team string) ([]e.GitlabConfig, error) {
	return g.repo.GetGitlabConfigsByTeam(ctx, team)
}
