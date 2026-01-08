package domain

type Technology struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Link        string   `json:"link"`
}

func (t *Technology) Fingerprint() string {
	return t.Name + "|" + t.Version + "|" + t.Description
}

func (t *Technology) Score() int {
	score := 0

	if t.Name != "" {
		score++
	}
	if t.Version != "" {
		score++
	}
	if t.Description != "" {
		score++
	}
	if t.Link != "" {
		score++
	}

	return score
}
