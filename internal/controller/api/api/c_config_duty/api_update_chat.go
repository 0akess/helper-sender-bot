package c_config_duty

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"helper-sender-bot/internal/entity"
	"net/http"
)

type updateChatReq struct {
	DutyTtlInMin       int    `json:"duty_ttl_in_minute" validate:"required"`
	DutyRepeatTtlInMin int    `json:"duty_repeat_ttl_in_minute" validate:"required"`
	EmojiStart         string `json:"emoji_start" validate:"required"`
	EmojiDone          string `json:"emoji_done" validate:"required"`
	WorkdayStart       int    `json:"workday_start" validate:"required"`
	WorkdayEnd         int    `json:"workday_end" validate:"required"`
}

func (c *CfgDutyController) updateChat(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return r.NotAuthMassage(err)
	}

	err = c.ucAuth.Auth(c.ctx, auth)
	if err != nil {
		return r.ForbiddenMassage(err)
	}

	channel := e.QueryParam("channel")
	if channel == "" {
		return r.InvalidInputMassage(fmt.Errorf("query 'channel' is required"))
	}

	var req updateChatReq
	if err := e.Bind(&req); err != nil {
		return r.InvalidInputMassage(err)
	}
	if err := e.Validate(&req); err != nil {
		return r.InvalidInputMassage(err)
	}

	chat := entity.Chat{
		DutyTtlInMinute:       req.DutyTtlInMin,
		DutyRepeatTtlInMinute: req.DutyRepeatTtlInMin,
		EmojiStart:            req.EmojiStart,
		EmojiDone:             req.EmojiDone,
		WorkdayStart:          req.WorkdayStart,
		WorkdayEnd:            req.WorkdayEnd,
	}

	if err := c.uc.UpdateDutyCfg(c.ctx, channel, auth.Team, chat); err != nil {
		return r.InternalErrorMassage(err)
	}
	return e.NoContent(http.StatusOK)
}
