package c_config_duty

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

func (c *CfgDutyController) deleteChat(e echo.Context) error {
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

	if err := c.dutyCfg.DeleteDutyCfg(e.Request().Context(), channel, auth.Team); err != nil {
		return responses.InternalErrorMessage(err)
	}
	return e.NoContent(http.StatusNoContent)
}
