package fetcher

import (
	"TechstackDetectorAPI/internal/core/domain"
	"context"
	"io"
	"net/url"
	"sync"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

type HTTPFetcher struct {
	client      *resty.Client
	maxParallel int
}

func NewHTTPFetcher(client *resty.Client, maxParallel int) *HTTPFetcher {
	return &HTTPFetcher{
		client:      client,
		maxParallel: maxParallel,
	}
}

func (f *HTTPFetcher) Fetch(
	ctx context.Context,
	plan domain.FetchPlan,
) map[string]*domain.HTTPResult {

	results := make(map[string]*domain.HTTPResult)
	execCache := make(map[string]*domain.HTTPResult)

	var mu sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(f.maxParallel)

	for _, r := range plan.Requests {
		r := r // capture loop variable

		g.Go(func() error {
			key := executionKey(plan.BaseURL, r)

			// cek cache
			mu.Lock()
			cached, ok := execCache[key]
			mu.Unlock()

			if ok {
				mu.Lock()
				results[r.ID] = cloneWithID(cached, r.ID)
				mu.Unlock()
				return nil
			}

			// execute request
			res := f.execute(ctx, plan.BaseURL, r)

			mu.Lock()
			execCache[key] = res
			results[r.ID] = cloneWithID(res, r.ID)
			mu.Unlock()

			return nil
		})
	}

	_ = g.Wait()
	return results
}

func (f *HTTPFetcher) execute(
	ctx context.Context,
	baseURL string,
	r domain.FetchRequest,
) *domain.HTTPResult {

	req := f.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true)

	target := joinURL(baseURL, r.Path)

	resp, err := req.Execute(r.HTTPMethod(), target)
	if err != nil {
		return &domain.HTTPResult{
			URL:   target,
			Error: err,
		}
	}
	defer resp.RawBody().Close()

	body, _ := io.ReadAll(resp.RawBody())

	finalURL := ""
	if resp.RawResponse != nil && resp.RawResponse.Request != nil {
		finalURL = resp.RawResponse.Request.URL.String()
	}

	return &domain.HTTPResult{
		URL:        target,
		FinalURL:   finalURL,
		StatusCode: resp.StatusCode(),
		Headers:    resp.Header(),
		Body:       body,
	}
}

func executionKey(baseURL string, r domain.FetchRequest) string {
	return r.HTTPMethod() + "|" + baseURL + "|" + r.Path
}

func cloneWithID(src *domain.HTTPResult, id string) *domain.HTTPResult {
	cp := *src
	cp.RequestID = id
	return &cp
}

func joinURL(base, path string) string {
	u, _ := url.JoinPath(base, path)
	return u
}
