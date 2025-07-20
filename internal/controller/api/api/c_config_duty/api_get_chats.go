package c_config_duty

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

type getChatsResponse struct {
	ChannelID          string `json:"channel_id"`
	DutyTtlInMin       int    `json:"duty_ttl_in_minute"`
	DutyRepeatTtlInMin int    `json:"duty_repeat_ttl_in_minute"`
	EmojiStart         string `json:"emoji_start"`
	EmojiDone          string `json:"emoji_done"`
	WorkdayStart       int    `json:"workday_start"`
	WorkdayEnd         int    `json:"workday_end"`
}

func (c *CfgDutyController) getChats(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return r.NotAuthMassage(err)
	}

	err = c.ucAuth.Auth(c.ctx, auth)
	if err != nil {
		return r.ForbiddenMassage(err)
	}

	chats, err := c.uc.GetListDutyCfgByTeam(c.ctx, auth.Team)
	if err != nil {
		return r.InternalErrorMassage(err)
	}

	if len(chats) == 0 {
		return e.JSON(http.StatusOK, []string{})
	}

	res := make([]getChatsResponse, len(chats))
	for i, chat := range chats {
		res[i] = getChatsResponse{
			ChannelID:          chat.ChannelID,
			DutyTtlInMin:       chat.DutyTtlInMinute,
			DutyRepeatTtlInMin: chat.DutyRepeatTtlInMinute,
			EmojiStart:         chat.EmojiStart,
			EmojiDone:          chat.EmojiDone,
			WorkdayStart:       chat.WorkdayStart,
			WorkdayEnd:         chat.WorkdayEnd,
		}
	}

	return e.JSON(http.StatusOK, res)
}
