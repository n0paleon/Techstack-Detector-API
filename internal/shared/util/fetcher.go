package util

import (
	"TechstackDetectorAPI/internal/core/domain"
	"net/url"
)

func ExecutionKey(baseURL string, r domain.FetchRequest) string {
	return r.HTTPMethod() + "|" + baseURL + "|" + r.Path
}

func CloneWithID(src *domain.HTTPResult, id string) *domain.HTTPResult {
	cp := *src
	cp.RequestID = id
	return &cp
}

func JoinURL(base, path string) string {
	u, _ := url.JoinPath(base, path)
	return u
}
