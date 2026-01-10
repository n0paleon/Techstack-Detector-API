package dto

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
