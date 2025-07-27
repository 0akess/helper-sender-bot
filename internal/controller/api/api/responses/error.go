package responses

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Status    int         `json:"-"`
	ErrDetail detailError `json:"error"`
}

type detailError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrDetail.Code, e.ErrDetail.Message)
}

func NewError(status int, code, message string, details interface{}) ResponseError {
	return ResponseError{
		Status: status,
		ErrDetail: detailError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

func (e ResponseError) Write(c echo.Context) {
	if err := c.JSON(e.Status, e); err != nil {
		c.Logger().Error("Failed to encode error response", "error", err)
	}
}

func BuildDetails(err error) map[string]string {
	if err == nil {
		return nil
	}

	var vErrs validator.ValidationErrors
	if !errors.As(err, &vErrs) {
		return map[string]string{"error": err.Error()}
	}

	details := make(map[string]string, len(vErrs))
	for _, fe := range vErrs {
		details[fe.Field()] = fe.Error()
	}
	return details
}

func InvalidInputMessage(err error) ResponseError {
	return NewError(
		http.StatusBadRequest,
		"INVALID_INPUT",
		"Add required field",
		BuildDetails(err),
	)
}

func NotAuthMessage(err error) ResponseError {
	return NewError(
		http.StatusUnauthorized,
		"NOT_AUTHORIZED",
		"Add auth data",
		BuildDetails(err),
	)
}

func ForbiddenMessage(err error) ResponseError {
	return NewError(
		http.StatusForbidden,
		"INVALID_AUTHORIZED_DATA",
		"Change auth data",
		BuildDetails(err),
	)
}

func InternalErrorMessage(err error) ResponseError {
	return NewError(
		http.StatusInternalServerError,
		"INTERNAL",
		"Internal Server Error",
		BuildDetails(err),
	)
}
