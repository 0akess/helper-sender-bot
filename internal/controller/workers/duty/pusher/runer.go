package pusher

import (
	"context"
	"time"
)

type uc interface {
	PusherBot(ctx context.Context)
}

type Pusher struct {
	uc uc
}

func NewPusher(uc uc) *Pusher {
	return &Pusher{
		uc: uc,
	}
}

func (p *Pusher) RunGoPusherBot(ctx context.Context, interval time.Duration) {
	go func() {
		for {
			p.uc.PusherBot(ctx)
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
			}
		}
	}()
}
