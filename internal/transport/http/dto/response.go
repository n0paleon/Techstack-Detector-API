package dto

import "TechstackDetectorAPI/internal/core/domain"

type Response struct {
	IsError bool   `json:"is_error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse(err error, data any) *Response {
	response := &Response{
		Data: data,
	}

	if err != nil {
		response.IsError = true
		response.Message = err.Error()
	}

	return response
}

type DetectorList struct {
	Name        string   `json:"name"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Link        string   `json:"link"`
}

func ConvertTechnologyToDetectorList(d domain.Technology) DetectorList {
	return DetectorList{
		Name:        d.Name,
		Tags:        d.Tags,
		Description: d.Description,
		Link:        d.Link,
	}
}
