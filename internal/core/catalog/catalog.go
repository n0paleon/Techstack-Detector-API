package catalog

import (
	"TechstackDetectorAPI/internal/core/domain"

	"github.com/labstack/gommon/log"
)

// NewTechnology generate technology based on a given detector ID. The detector ID must be the same as the value returned by the Name() method on ports.Detector
func NewTechnology(id DetectorID, version string) domain.Technology {
	meta, ok := registry[id]
	if !ok {
		log.Errorf("Technology %s not found in registry", id)
		return domain.Technology{}
	}

	return domain.Technology{
		Name:        meta.Name,
		Version:     version,
		Tags:        meta.Tags,
		Description: meta.Description,
		Link:        meta.Link,
	}
}

// List return all available technology
func List() []domain.Technology {
	results := make([]domain.Technology, 0, len(registry))
	for _, tech := range registry {
		results = append(results, tech)
	}

	return results
}
