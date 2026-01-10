package dto

type FastDetectRequest struct {
	Target string `json:"target" validate:"required,url,lte=500"` // max 500 chars
}
