package ports

import (
	"TechstackDetectorAPI/internal/core/domain"
	"context"
)

// Fetcher menentukan kontrak bagaimana fetcher harus bekerja (HTTP? WS? Chromium?)
type Fetcher interface {
	Fetch(ctx context.Context, plan *domain.FetchPlan) (*domain.FetchContext, error)
}
