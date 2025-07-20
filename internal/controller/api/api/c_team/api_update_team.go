package c_chat_config

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/controller/api/api/middleware"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"helper-sender-bot/internal/entity"
	"net/http"
)

type updateTeamReq struct {
	NewToken   uuid.UUID `json:"new_token" validate:"required"`
	NewLeadEID string    `json:"new_team_lead_eid" validate:"required,min=1"`
}

func (c *Controller) updateTeam(e echo.Context) error {
	auth, err := middleware.GetAuth(e)
	if err != nil {
		return r.NotAuthMassage(err)
	}

	err = c.ucAuth.Auth(c.ctx, auth)
	if err != nil {
		return r.ForbiddenMassage(err)
	}
	var req updateTeamReq
	if err := e.Bind(&req); err != nil {
		return r.InvalidInputMassage(err)
	}
	if err := e.Validate(req); err != nil {
		return r.InvalidInputMassage(err)
	}

	newTeam := entity.Team{
		Token:   req.NewToken,
		LeadEID: req.NewLeadEID,
	}

	err = c.uc.UpdateTeam(c.ctx, newTeam, auth.Team, auth.Token)
	if err != nil {
		return r.InternalErrorMassage(err)
	}
	return e.JSON(http.StatusOK, map[string]bool{"success": true})
}
