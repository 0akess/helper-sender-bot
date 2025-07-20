package c_chat_config

import (
	"github.com/labstack/echo/v4"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"net/http"
)

type getAllTeamQuery struct {
	Cursor         int    `query:"cursor"`
	Limit          int    `query:"limit"`
	SearchTeamName string `query:"team_name_like"`
}

func (c *Controller) getTeam(e echo.Context) error {
	var query getAllTeamQuery
	if err := e.Bind(&query); err != nil {
		return r.InvalidInputMassage(err)
	}
	if err := e.Validate(query); err != nil {
		return r.InvalidInputMassage(err)
	}

	if query.Cursor <= 0 {
		query.Cursor = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}

	teams, nextCursor, err := c.uc.GetTeams(c.ctx, query.Limit, query.Cursor-1, query.SearchTeamName)
	if err != nil {
		return r.InternalErrorMassage(err)
	}

	if len(teams) == 0 {
		return e.JSON(http.StatusOK, []string{})
	}

	response := map[string]interface{}{
		"next_cursor": nextCursor,
		"teams":       teams,
	}
	return e.JSON(http.StatusOK, response)
}
