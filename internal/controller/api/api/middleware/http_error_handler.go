package middleware

import (
	"errors"
	"helper-sender-bot/internal/controller/api/api/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var resp responses.ResponseError
		if errors.As(err, &resp) {
			resp.Write(c)
			return
		}

		var ehe *echo.HTTPError
		if errors.As(err, &ehe) {
			code, msg := normalizeEchoError(ehe)
			writeJSON(c, code, map[string]any{
				"error": map[string]any{
					"code":    http.StatusText(code),
					"message": msg,
				},
			})
			return
		}

		c.Logger().Error("HTTP Error",
			"err", err,
			"path", c.Path(),
			"method", c.Request().Method,
			"request_id", c.Response().Header().Get("X-Request-Id"),
		)

		writeJSON(c, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": http.StatusText(http.StatusInternalServerError),
			},
		})
	}
}

func writeJSON(c echo.Context, code int, data interface{}) {
	if err := c.JSON(code, data); err != nil {
		c.Logger().Error("Failed to encode error response", "error", err)
	}
}

func normalizeEchoError(he *echo.HTTPError) (int, string) {
	code := he.Code
	msg, _ := he.Message.(string)

	switch {
	case code == http.StatusServiceUnavailable:
		code = http.StatusGatewayTimeout
		msg = "Gateway Timeout"
	case msg == "Unauthorized":
		code = http.StatusForbidden
		msg = "Token is invalid"
	}
	return code, msg
}
