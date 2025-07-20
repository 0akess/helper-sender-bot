package c_config_duty

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"helper-sender-bot/internal/entity"
	"net/http"
)

type createChatReq struct {
	ChannelID             string `json:"channel_id" validate:"required"`
	DutyTtlInMinute       int    `json:"duty_ttl_in_minute" validate:"required,min=1"`
	DutyRepeatTtlInMinute int    `json:"duty_repeat_ttl_in_minute" validate:"required,min=1"`
	EmojiStart            string `json:"emoji_start" validate:"required"`
	EmojiDone             string `json:"emoji_done" validate:"required"`
	WorkdayStart          int    `json:"workday_start" validate:"required"`
	WorkdayEnd            int    `json:"workday_end" validate:"required"`
}

func (c *CfgDutyController) createChat(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return r.NotAuthMassage(err)
	}

	err = c.ucAuth.Auth(c.ctx, auth)
	if err != nil {
		return r.ForbiddenMassage(err)
	}

	var req createChatReq
	if err := e.Bind(&req); err != nil {
		return r.InvalidInputMassage(err)
	}

	if err := e.Validate(req); err != nil {
		return r.InvalidInputMassage(err)
	}

	chat := entity.Chat{
		ChannelID:             req.ChannelID,
		DutyTtlInMinute:       req.DutyTtlInMinute,
		DutyRepeatTtlInMinute: req.DutyRepeatTtlInMinute,
		EmojiStart:            req.EmojiStart,
		EmojiDone:             req.EmojiDone,
		WorkdayStart:          req.WorkdayStart,
		WorkdayEnd:            req.WorkdayEnd,
	}

	err = c.uc.CreateDutyCfg(c.ctx, chat, auth.Team)
	if err != nil {
		return r.InternalErrorMassage(err)
	}
	return e.JSON(http.StatusCreated, map[string]bool{"success": true})
}
