package middleware

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/entity"
	"net/http"
)

type ctxKey string

const (
	ctxKeyTeam  ctxKey = "team"
	ctxKeyToken ctxKey = "token"
)

func HeaderAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		team := c.Request().Header.Get("X-Team")
		if team == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "'X-Team' header обязательный")
		}
		tokRaw := c.Request().Header.Get("X-Auth-Token")
		token, err := uuid.Parse(tokRaw)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "'X-Auth-Token' должен быть UUID")
		}

		c.Set(string(ctxKeyTeam), team)
		c.Set(string(ctxKeyToken), token)

		return next(c)
	}
}

func GetAuth(e echo.Context) (auth entity.AuthMeta, err error) {
	vToken, okToken := e.Get(string(ctxKeyToken)).(uuid.UUID)
	vTeam, okTeam := e.Get(string(ctxKeyTeam)).(string)
	if !okToken || !okTeam {
		return entity.AuthMeta{}, fmt.Errorf("в контексте нет всех данных авторизации")
	}
	return entity.AuthMeta{
		Team:  vTeam,
		Token: vToken,
	}, nil
}
