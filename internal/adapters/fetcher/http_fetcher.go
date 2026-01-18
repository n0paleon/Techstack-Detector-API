package fetcher

import (
	"TechstackDetectorAPI/internal/core/domain"
	"TechstackDetectorAPI/internal/shared/util"
	"context"
	"io"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"

	fakeua "github.com/EDDYCJY/fake-useragent"
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
			key := util.ExecutionKey(plan.BaseURL, r)

			// cek cache
			mu.Lock()
			cached, ok := execCache[key]
			mu.Unlock()

			if ok {
				mu.Lock()
				results[r.ID] = util.CloneWithID(cached, r.ID)
				mu.Unlock()
				return nil
			}

			// execute request
			res := f.execute(ctx, plan.BaseURL, r)
			if res == nil {
				return nil
			}

			mu.Lock()
			execCache[key] = res
			results[r.ID] = util.CloneWithID(res, r.ID)
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
	ua := fakeua.Mobile()

	req := f.client.R().
		SetContext(ctx).
		SetHeader("User-Agent", ua).
		SetDoNotParseResponse(true)

	target := util.JoinURL(baseURL, r.Path)

	resp, err := req.Execute(r.HTTPMethod(), target)
	if err != nil {
		if resp == nil || resp.StatusCode() == 0 {
			log.Errorf("http_fetcher return nil response, err: %v", err)
			return nil
		}

		return &domain.HTTPResult{
			RequestID:  r.ID,
			StatusCode: resp.StatusCode(),
			Headers:    resp.Header(),
			Body:       resp.Body(),
			URL:        target,
			Error:      err,
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
