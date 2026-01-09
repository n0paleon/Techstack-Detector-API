package config

import (
	"TechstackDetectorAPI/internal/core/domain"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	IsError bool   `json:"is_error"`
	Message string `json:"message"`
}

func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	var echoError *echo.HTTPError
	switch {
	case errors.As(err, &echoError):
		code = echoError.Code
		message = echoError.Message.(string)
	case errors.Is(err, domain.ErrInvalidTarget):
		code = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, domain.ErrBlockedTarget):
		code = http.StatusForbidden
		message = err.Error()
	}

	errorJSON := ErrorResponse{
		IsError: true,
		Message: message,
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, errorJSON)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
