package updater_posts

import (
	"context"
	"time"
)

type uc interface {
	UpdaterPosts(ctx context.Context, intervalCycle time.Duration)
}

type UpdaterPosts struct {
	uc uc
}

func NewUpdaterPosts(uc uc) *UpdaterPosts {
	return &UpdaterPosts{
		uc: uc,
	}
}

func (up *UpdaterPosts) RunGoUpdaterPosts(ctx context.Context, intervalCycle, intervalGorRun time.Duration) {
	go func() {
		for {
			up.uc.UpdaterPosts(ctx, intervalCycle)
			select {
			case <-ctx.Done():
				return
			case <-time.After(intervalGorRun):
			}
		}
	}()
}
