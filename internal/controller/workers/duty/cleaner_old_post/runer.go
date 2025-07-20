package cleaner_old_post

import (
	"context"
	"time"
)

type uc interface {
	CleanerOldPost(ctx context.Context, intervalCycle time.Duration)
}

type Cleaner struct {
	uc uc
}

func NewCleaner(uc uc) *Cleaner {
	return &Cleaner{
		uc: uc,
	}
}

func (c *Cleaner) RunGoCleanerOldPost(ctx context.Context, intervalCycle, intervalGoRun time.Duration) {
	go func() {
		for {
			c.uc.CleanerOldPost(ctx, intervalCycle)
			select {
			case <-ctx.Done():
				return
			case <-time.After(intervalGoRun):
			}
		}
	}()
}
