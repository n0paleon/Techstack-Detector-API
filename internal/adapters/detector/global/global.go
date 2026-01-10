package global

import (
	"TechstackDetectorAPI/internal/core/domain"

	"github.com/brianvoe/gofakeit/v7"
)

type Detector struct {
}

var basePlan = []domain.FetchRequest{
	{
		ID:          "root",
		Path:        "/",
		Method:      "GET",
		Description: "default homepage",
	},
}

// please do not edit anything on this file except the FetchPlan
// for more info, read README.md

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) Name() string {
	return "GLOBAL_DETECTOR"
}

// FetchPlan in this file is used to manage a list of common requests such as homepage, random 404 requests, etc.
// for more info, please read the README.md
func (d *Detector) FetchPlan(baseUrl string) *domain.FetchPlan {
	reqs := make([]domain.FetchRequest, 0, len(basePlan)+1)
	reqs = append(reqs, basePlan...)

	reqs = append(reqs, domain.FetchRequest{
		ID:          "random-404-error",
		Path:        gofakeit.LetterN(25),
		Method:      "GET",
		Description: "trigger random 404 error",
	})

	return &domain.FetchPlan{
		BaseURL:  baseUrl,
		Requests: reqs,
	}
}

func (d *Detector) Detect(_ *domain.FetchContext) ([]domain.Technology, error) {
	return []domain.Technology{}, nil
}
