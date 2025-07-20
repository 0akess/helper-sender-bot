package c_config_gitlab

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

func (c *Controller) deleteGitCfg(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return r.NotAuthMassage(err)
	}

	err = c.ucAuth.Auth(c.ctx, auth)
	if err != nil {
		return r.ForbiddenMassage(err)
	}

	gitURL, projectID, err := gitUrlAndIdQuery(e)
	if err != nil {
		return r.InvalidInputMassage(err)
	}

	if err = c.uc.DeleteGitlabConfig(c.ctx, projectID, gitURL, auth.Team); err != nil {
		return r.InternalErrorMassage(err)
	}
	return e.NoContent(http.StatusNoContent)
}
