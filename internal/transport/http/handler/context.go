package handler

import (
	"context"

	"github.com/labstack/echo/v4"
)

func getRequestContext(e echo.Context) context.Context {
	return e.Request().Context()
}
