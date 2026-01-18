package global

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"TechstackDetectorAPI/internal/core/domain"

	"github.com/brianvoe/gofakeit/v7"
)

type Detector struct{}

var basePlan = []domain.FetchRequest{
	{
		ID:          "root",
		Path:        "/",
		Method:      "GET",
		Description: "default homepage",
	},
	{
		ID:          "global-env",
		Path:        "/.env",
		Method:      "GET",
		Description: "Check for .env file exposure",
	},
}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) ID() catalog.DetectorID {
	return catalog.GLOBAL
}

func (d *Detector) FetchPlan(baseURL string) *domain.FetchPlan {
	// Create requests by copying basePlan and adding the random 404 request
	requests := make([]domain.FetchRequest, len(basePlan), len(basePlan)+1)
	copy(requests, basePlan)

	requests = append(requests,
		domain.FetchRequest{
			ID:          "random-404-error",
			Path:        gofakeit.LetterN(25),
			Method:      "GET",
			Description: "trigger random 404 error",
		},
		domain.FetchRequest{
			ID:          "random-404-error-post-method",
			Path:        gofakeit.LetterN(25),
			Method:      "POST",
			Description: "trigger random 404 error with POST method",
		},
	)

	return &domain.FetchPlan{
		BaseURL:  baseURL,
		Requests: requests,
	}
}

func (d *Detector) Detect(_ *domain.FetchContext) ([]domain.Technology, error) {
	return nil, nil
}
