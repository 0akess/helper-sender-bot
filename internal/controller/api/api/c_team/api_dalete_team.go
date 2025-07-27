package c_chat_config

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

func (t *Controller) deleteTeam(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return responses.NotAuthMessage(err)
	}

	err = t.auth.CheckAuth(e.Request().Context(), auth)
	if err != nil {
		return responses.ForbiddenMessage(err)
	}

	err = t.team.DeleteTeam(e.Request().Context(), auth.Team, auth.Token)
	if err != nil {
		return responses.InternalErrorMessage(err)
	}
	return e.JSON(http.StatusNoContent, map[string]bool{"success": true})
}
