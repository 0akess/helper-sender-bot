package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"helper-sender-bot/internal/entity"
)

type repo interface {
	OkTokenTeam(ctx context.Context, teamName string, token uuid.UUID) (bool, error)
}

type Auth struct {
	ctx  context.Context
	repo repo
}

func NewAuth(ctx context.Context, repo repo) *Auth {
	return &Auth{
		ctx:  ctx,
		repo: repo,
	}
}

func (a *Auth) Auth(ctx context.Context, auth entity.AuthMeta) error {
	tokenOk, err := a.repo.OkTokenTeam(ctx, auth.Team, auth.Token)
	if err != nil {
		return err
	}
	if !tokenOk {
		return fmt.Errorf("invalid token or team")
	}
	return nil
}
