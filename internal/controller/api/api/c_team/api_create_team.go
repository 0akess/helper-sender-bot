package c_chat_config

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	r "helper-sender-bot/internal/controller/api/api/responses"
	"helper-sender-bot/internal/entity"
	"net/http"
)

type teamReq struct {
	Name    string    `json:"team_name" validate:"required,min=2"`
	Token   uuid.UUID `json:"token" validate:"required"`
	LeadEID string    `json:"team_lead_eid" validate:"required,min=1"`
}

func (c *Controller) createTeam(e echo.Context) error {
	var req teamReq
	if err := e.Bind(&req); err != nil {
		return r.InvalidInputMassage(err)
	}
	if err := e.Validate(req); err != nil {
		return r.InvalidInputMassage(err)
	}

	team := entity.Team{
		Name:    req.Name,
		LeadEID: req.LeadEID,
		Token:   req.Token,
	}

	err := c.uc.CreateTeam(c.ctx, team)
	if err != nil {
		return r.InternalErrorMassage(err)
	}
	return e.JSON(http.StatusCreated, map[string]bool{"success": true})
}
