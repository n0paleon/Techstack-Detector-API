package config

import (
	"TechstackDetectorAPI/internal/core/domain"
	"TechstackDetectorAPI/internal/transport/http/dto"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

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

	errorJSON := dto.NewResponse(errors.New(message), nil)

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
