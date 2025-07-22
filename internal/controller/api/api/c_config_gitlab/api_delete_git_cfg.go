package c_config_gitlab

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	"helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

func (c *Controller) deleteGitCfg(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return responses.NotAuthMassage(err)
	}

	err = c.auth.CheckAuth(e.Request().Context(), auth)
	if err != nil {
		return responses.ForbiddenMassage(err)
	}

	gitURL, projectID, err := gitUrlAndIdQuery(e)
	if err != nil {
		return responses.InvalidInputMassage(err)
	}

	if err = c.gitlabCfg.DeleteGitlabConfig(e.Request().Context(), projectID, gitURL, auth.Team); err != nil {
		return responses.InternalErrorMassage(err)
	}
	return e.NoContent(http.StatusNoContent)
}
