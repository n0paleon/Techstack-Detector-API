package ports

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"TechstackDetectorAPI/internal/core/domain"
)

type Detector interface {
	// ID harus mereturn ID unik untuk tiap detector
	ID() catalog.DetectorID
	// Detect harus bisa mendeteksi technology apa yang digunakan berdasarkan domain.FetchContext yang diberikan
	Detect(ctx *domain.FetchContext) ([]domain.Technology, error)
	FetchPlan(target string) *domain.FetchPlan // optional, boleh nil
}
