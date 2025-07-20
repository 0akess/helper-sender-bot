package pingonsla

import (
	"context"
	"time"
)

type uc interface {
	SendRepeatPush(ctx context.Context)
}

type RepeatPush struct {
	uc uc
}

func NewRepeatPush(uc uc) *RepeatPush {
	return &RepeatPush{
		uc: uc,
	}
}

func (rp *RepeatPush) RunGoSendRepeatPush(ctx context.Context, interval time.Duration) {
	go func() {
		for {
			rp.uc.SendRepeatPush(ctx)
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
			}
		}
	}()
}
