package echohttp

import (
	"TechstackDetectorAPI/internal/transport/http/echo/config"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer() *echo.Echo {
	e := echo.New()
	e.JSONSerializer = &config.SonicJSONSerializer{}
	e.Server.WriteTimeout = 30 * time.Second // 30s max write timeout
	e.HTTPErrorHandler = config.HTTPErrorHandler

	// middleware
	e.Use(middleware.Recover())

	// TODO: implement advance server setup

	return e
}
