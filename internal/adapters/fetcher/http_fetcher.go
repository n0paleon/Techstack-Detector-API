package fetcher

import (
	"context"
	"io"
	"net/url"
	"sync"

	"TechstackDetectorAPI/internal/core/domain"

	"github.com/go-resty/resty/v2"
)

type HTTPFetcher struct {
	client *resty.Client
}

func NewHTTPFetcher(client *resty.Client) *HTTPFetcher {
	return &HTTPFetcher{client: client}
}

func (f *HTTPFetcher) Fetch(
	ctx context.Context,
	plan domain.FetchPlan,
) map[string]*domain.HTTPResult {

	results := make(map[string]*domain.HTTPResult)
	execCache := make(map[string]*domain.HTTPResult)

	var mu sync.Mutex

	for _, r := range plan.Requests {
		key := executionKey(plan.BaseURL, r)

		mu.Lock()
		cached, ok := execCache[key]
		mu.Unlock()

		if ok {
			results[r.ID] = cloneWithID(cached, r.ID)
			continue
		}

		res := f.execute(ctx, plan.BaseURL, r)

		mu.Lock()
		execCache[key] = res
		results[r.ID] = cloneWithID(res, r.ID)
		mu.Unlock()
	}

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
