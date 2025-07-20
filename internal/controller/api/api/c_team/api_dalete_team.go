package c_chat_config

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

func (c *Controller) deleteTeam(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return r.NotAuthMassage(err)
	}

	err = c.ucAuth.Auth(c.ctx, auth)
	if err != nil {
		return r.ForbiddenMassage(err)
	}

	err = c.uc.DeleteTeam(c.ctx, auth.Team, auth.Token)
	if err != nil {
		return r.InternalErrorMassage(err)
	}
	return e.JSON(http.StatusNoContent, map[string]bool{"success": true})
}
