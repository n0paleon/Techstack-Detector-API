package webserver

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"TechstackDetectorAPI/internal/core/domain"
	"strings"
)

type LiteSpeedDetector struct {
}

func NewLiteSpeed() *LiteSpeedDetector {
	return &LiteSpeedDetector{}
}

func (d *LiteSpeedDetector) ID() catalog.DetectorID {
	return catalog.LiteSpeed
}

func (d *LiteSpeedDetector) FetchPlan(_ string) *domain.FetchPlan {
	return nil
}

func (d *LiteSpeedDetector) Detect(ctx *domain.FetchContext) ([]domain.Technology, error) {
	score := 0

	for _, httpResult := range ctx.HTTP {
		if httpResult == nil || httpResult.Error != nil {
			continue
		}

		header := httpResult.Headers.Get("Server")
		if strings.Contains(strings.ToLower(header), "litespeed") {
			score++
		}

		if httpResult.Headers.Get("X-Turbo-Charged-By") != "" {
			score++
		}

		if strings.Contains(string(httpResult.Body), "Proudly powered by LiteSpeed Web Server") {
			score++
		}

		if httpResult.Headers.Get("X-LiteSpeed-Cache") != "" || httpResult.Headers.Get("X-LiteSpeed-Purge") != "" {
			score++
		}
	}

	var results []domain.Technology
	if score >= 1 {
		results = append(results, catalog.NewTechnology(d.ID(), ""))
	}

	return results, nil
}
