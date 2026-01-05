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
