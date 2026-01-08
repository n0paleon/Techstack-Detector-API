package php

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"TechstackDetectorAPI/internal/core/domain"
	"strings"
)

type PHPDetector struct{}

func NewPHPDetector() *PHPDetector {
	return &PHPDetector{}
}

func (d *PHPDetector) Name() string {
	return "php"
}

func (d *PHPDetector) FetchPlan(target string) *domain.FetchPlan {
	return &domain.FetchPlan{
		BaseURL: target,
		Requests: []domain.FetchRequest{
			{
				ID:          "php-probing",
				Path:        "/index.php",
				Method:      "GET",
				Description: "Probing for PHP indicators",
			},
		},
	}
}

func (d *PHPDetector) Detect(ctx *domain.FetchContext) ([]domain.Technology, error) {
	detected := false
	version := ""

	for _, res := range ctx.HTTP {
		if res.Error != nil {
			continue
		}

		if val := res.Headers.Get("X-Powered-By"); val != "" {
			lowerVal := strings.ToLower(val)
			if strings.Contains(lowerVal, "php") {
				detected = true
				parts := strings.Split(val, "/")
				if len(parts) > 1 {
					version = parts[1]
				}
			}
		}

		for _, cookie := range res.Headers.Values("Set-Cookie") {
			if strings.Contains(cookie, "PHPSESSID") {
				detected = true
				break
			}
		}

		if res.Headers.Get("X-PHP-Originating-Script") != "" {
			detected = true
		}

		bodyStr := string(res.Body)
		if strings.Contains(bodyStr, "fatal error") && strings.Contains(bodyStr, ".php on line") {
			detected = true
		}

		if strings.Contains(bodyStr, ".php\"") || strings.Contains(bodyStr, ".php'") {
			detected = true
		}

		if detected {
			break
		}
	}

	if detected {
		tech := catalog.PHP("")
		if version != "" {
			tech.Version = version
		}
		return []domain.Technology{*tech}, nil
	}

	return nil, nil
}
