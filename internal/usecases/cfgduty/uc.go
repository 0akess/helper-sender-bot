package cfgduty

import (
	"context"
	"helper-sender-bot/internal/entity"
)

type repo interface {
	CreateCfgDuty(ctx context.Context, c entity.Chat, team string) error
	GetListCfgDutyByTeam(ctx context.Context, team string) ([]entity.Chat, error)
	UpdateCfgDuty(ctx context.Context, team, channel string, upd entity.Chat) error
	DeleteCfgDuty(ctx context.Context, team, channel string) error
}

type DutyCfg struct {
	ctx  context.Context
	repo repo
}

func NewDutyCfgCases(ctx context.Context, repo repo) *DutyCfg {
	return &DutyCfg{
		ctx:  ctx,
		repo: repo,
	}
}

func (cc *DutyCfg) GetListDutyCfgByTeam(ctx context.Context, team string) ([]entity.Chat, error) {
	return cc.repo.GetListCfgDutyByTeam(ctx, team)
}

func (cc *DutyCfg) CreateDutyCfg(ctx context.Context, chat entity.Chat, team string) error {
	return cc.repo.CreateCfgDuty(ctx, chat, team)
}

func (cc *DutyCfg) UpdateDutyCfg(ctx context.Context, channel, team string, upd entity.Chat) error {
	return cc.repo.UpdateCfgDuty(ctx, team, channel, upd)
}

func (cc *DutyCfg) DeleteDutyCfg(ctx context.Context, channel, team string) error {
	return cc.repo.DeleteCfgDuty(ctx, team, channel)
}
