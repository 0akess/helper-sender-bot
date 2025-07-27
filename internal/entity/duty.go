package entity

import "time"

type Chat struct {
	Team                  string
	ChannelID             string
	DutyTtlInMinute       int
	DutyRepeatTtlInMinute int
	EmojiStart            string
	EmojiDone             string
	WorkdayStart          int
	WorkdayEnd            int
}

type PostsInfoDuty struct {
	ChannelID  string
	PostID     string
	CreateAt   time.Time
	LastPushAt time.Time
	InProgress bool
}

// IsNotWorkingHours фильтр определяет в какие часы нельзя слать уведомления
func (c *Chat) IsNotWorkingHours() bool {
	msk := time.FixedZone("MSK", 3*60*60)
	h := time.Now().In(msk).Hour()
	return h < c.WorkdayStart || h >= c.WorkdayEnd
}
