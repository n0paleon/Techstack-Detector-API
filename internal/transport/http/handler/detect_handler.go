package handler

import (
	"TechstackDetectorAPI/internal/core/service"
	"TechstackDetectorAPI/internal/transport/http/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DetectHandler struct {
	svc *service.DetectionService
}

// TODO: add url filtering before processing

func (h *DetectHandler) DetectorList(c echo.Context) error {
	list := h.svc.DetectorList()
	results := make([]dto.DetectorList, 0, len(list))

	for _, d := range list {
		results = append(results, dto.ConvertTechnologyToDetectorList(d))
	}

	return c.JSON(http.StatusOK, dto.NewResponse(nil, echo.Map{
		"detectors": results,
	}))
}

func (h *DetectHandler) FastDetect(c echo.Context) error {
	request := new(dto.FastDetectRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// validate json body
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewResponse(err, nil))
	}

	ctx := getRequestContext(c)
	data, err := h.svc.Detect(ctx, request.Target)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.NewResponse(err, nil))
	}

	return c.JSON(http.StatusOK, dto.NewResponse(nil, echo.Map{
		"result": data,
	}))
}

func NewDetectHandler(svc *service.DetectionService) *DetectHandler {
	return &DetectHandler{
		svc: svc,
	}
}
