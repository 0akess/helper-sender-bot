package ping_on_sla

import (
	"context"
	"time"
)

type uc interface {
	SendDayPing(ctx context.Context)
}

type DayPinger struct {
	uc uc
}

func NewDayPinger(uc uc) *DayPinger {
	return &DayPinger{
		uc: uc,
	}
}

func (dp *DayPinger) RunGoSendDayPinger(ctx context.Context) {
	go func() {
		for {
			msk := time.FixedZone("MSK", 3*60*60)
			now := time.Now().In(msk)

			next := time.Date(
				now.Year(), now.Month(), now.Day(),
				14, 0, 0, 0, msk,
			)
			if !next.After(now) {
				next = next.Add(24 * time.Hour)
			}

			time.Sleep(time.Until(next))

			wd := next.Weekday()
			if wd >= time.Monday && wd <= time.Friday {
				dp.uc.SendDayPing(ctx)
			}
		}
	}()
}
