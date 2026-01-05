package ports

import (
	"TechstackDetectorAPI/internal/core/domain"
)

type Detector interface {
	// Name harus mereturn nama unik yang merepresentasikan adapter tersebut bekerja untuk technology apa
	Name() string
	// Detect harus bisa mendeteksi technology apa yang digunakan berdasarkan domain.FetchContext yang diberikan
	Detect(ctx *domain.FetchContext) ([]domain.Technology, error)
	FetchPlan(target string) *domain.FetchPlan // optional, boleh nil
}
