package catalog

import "TechstackDetectorAPI/internal/core/domain"

func PHP(version string) *domain.Technology {
	return &domain.Technology{
		Name:        "PHP",
		Version:     version,
		Tags:        []string{"programming language", "php"},
		Description: "PHP (recursive acronym for PHP: Hypertext Preprocessor) is a widely-used open source general-purpose scripting language that is especially suited for web development and can be embedded into HTML.",
		Link:        "https://www.php.net/manual/en/introduction.php",
	}
}
