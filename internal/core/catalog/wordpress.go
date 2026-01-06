package catalog

import "TechstackDetectorAPI/internal/core/domain"

func WordPress(version string) *domain.Technology {
	return &domain.Technology{
		Name:        "WordPress",
		Tags:        []string{"cms", "php", "website"},
		Version:     version,
		Description: "WordPress is an open-source CMS.",
		Link:        "https://wordpress.org/",
	}
}
