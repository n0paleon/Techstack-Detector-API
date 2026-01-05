package echohttp

import (
	"TechstackDetectorAPI/internal/transport/http/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	detectHandler *handler.DetectHandler,
) {
	v1 := e.Group("/v1")

	v1.GET("/fast-detect", detectHandler.FastDetect)
}
