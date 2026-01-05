package handler

import (
	"TechstackDetectorAPI/internal/core/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DetectHandler struct {
	svc *service.DetectionService
}

func (h *DetectHandler) FastDetect(c echo.Context) error {
	target := c.QueryParam("target")
	if target == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"is_error": true,
			"message":  "invalid target url",
		})
	}

	ctx := getRequestContext(c)
	data, err := h.svc.Detect(ctx, target)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"is_error": true,
			"message":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"is_error": false,
		"data":     data,
	})
}

func NewDetectHandler(svc *service.DetectionService) *DetectHandler {
	return &DetectHandler{
		svc: svc,
	}
}
