package c_config_duty

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

func (c *CfgDutyController) deleteChat(e echo.Context) error {
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

	if err := c.uc.DeleteDutyCfg(c.ctx, channel, auth.Team); err != nil {
		return r.InternalErrorMassage(err)
	}
	return e.NoContent(http.StatusNoContent)
}
