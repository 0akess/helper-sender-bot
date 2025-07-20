package team

import (
	"context"
	"github.com/google/uuid"
	"helper-sender-bot/internal/entity"
)

type repo interface {
	CreateTeam(ctx context.Context, team entity.Team) error
	GetListTeam(ctx context.Context, limit, cursor int, search string) ([]string, int, error)
	UpdateTeam(ctx context.Context, teamName string, token uuid.UUID, newTeam entity.Team) error
	DeleteTeam(ctx context.Context, teamName string, token uuid.UUID) error
}

type Team struct {
	ctx  context.Context
	repo repo
}

func NewTeamCases(ctx context.Context, repo repo) *Team {
	return &Team{
		ctx:  ctx,
		repo: repo,
	}
}

func (t *Team) CreateTeam(ctx context.Context, team entity.Team) error {
	return t.repo.CreateTeam(ctx, team)
}

func (t *Team) GetTeams(ctx context.Context, limit, cursor int, search string) ([]string, int, error) {
	return t.repo.GetListTeam(ctx, limit, cursor, search)
}

func (t *Team) UpdateTeam(ctx context.Context, newTeam entity.Team, teamName string, token uuid.UUID) error {
	return t.repo.UpdateTeam(ctx, teamName, token, newTeam)
}

func (t *Team) DeleteTeam(ctx context.Context, teamName string, token uuid.UUID) error {
	return t.repo.DeleteTeam(ctx, teamName, token)
}
