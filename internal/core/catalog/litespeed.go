package catalog

import "TechstackDetectorAPI/internal/core/domain"

func LiteSpeed(version string) *domain.Technology {
	return &domain.Technology{
		Name:        "LiteSpeed",
		Version:     version,
		Tags:        []string{"web server", "proxy", "proprietary software"},
		Description: "LiteSpeed is a high-performance, high-scalability web server with comprehensive features.",
		Link:        "https://en.wikipedia.org/wiki/LiteSpeed_Web_Server",
	}
}
