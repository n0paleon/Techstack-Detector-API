package domain

type Technology struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Link        string   `json:"link"`
}

// Fingerprint returns a unique identifier for the technology.
// This is used to deduplicate technologies with the same name.
func (t *Technology) Fingerprint() string {
	// For deduplication, we primarily care about name and version
	// Description and link might differ between detectors, but they
	// should refer to the same technology
	return t.Name + "|" + t.Version
}

func (t *Technology) Score() int {
	score := 0

	if t.Name != "" {
		score += 10 // Base score for detection
	}
	if t.Version != "" {
		score += 5 // Bonus for having version
	}
	if t.Description != "" {
		score += 2 // Small bonus for description
	}
	if t.Link != "" {
		score += 1 // Small bonus for link
	}

	return score
}
