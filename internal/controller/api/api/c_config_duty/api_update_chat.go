package c_config_duty

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/controller/api/api/responses"
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
		return responses.NotAuthMessage(err)
	}

	err = c.auth.CheckAuth(e.Request().Context(), auth)
	if err != nil {
		return responses.ForbiddenMessage(err)
	}

	channel := e.QueryParam("channel")
	if channel == "" {
		return responses.InvalidInputMessage(fmt.Errorf("query 'channel' is required"))
	}

	var req updateChatReq
	if err := e.Bind(&req); err != nil {
		return responses.InvalidInputMessage(err)
	}
	if err := e.Validate(&req); err != nil {
		return responses.InvalidInputMessage(err)
	}

	chat := entity.Chat{
		DutyTtlInMinute:       req.DutyTtlInMin,
		DutyRepeatTtlInMinute: req.DutyRepeatTtlInMin,
		EmojiStart:            req.EmojiStart,
		EmojiDone:             req.EmojiDone,
		WorkdayStart:          req.WorkdayStart,
		WorkdayEnd:            req.WorkdayEnd,
	}

	if err := c.dutyCfg.UpdateDutyCfg(e.Request().Context(), channel, auth.Team, chat); err != nil {
		return responses.InternalErrorMessage(err)
	}
	return e.NoContent(http.StatusOK)
}
