package fetcher

import (
	"TechstackDetectorAPI/internal/core/domain"
	"context"
	"errors"
	"net/url"

	"golang.org/x/sync/errgroup"
)

type Orchestrator struct {
	httpFetcher *HTTPFetcher
	dnsFetcher  *DNSFetcher
}

func New(
	httpFetcher *HTTPFetcher,
	dnsFetcher *DNSFetcher,
) *Orchestrator {
	return &Orchestrator{
		httpFetcher: httpFetcher,
		dnsFetcher:  dnsFetcher,
	}
}

func (o *Orchestrator) Fetch(
	ctx context.Context,
	plan *domain.FetchPlan,
) (*domain.FetchContext, error) {

	u, err := url.Parse(plan.BaseURL)
	if err != nil {
		return nil, err
	}

	g, gCtx := errgroup.WithContext(ctx)

	var dnsResult *domain.DNSResult
	httpResults := make(map[string]*domain.HTTPResult)

	// dns fetch
	g.Go(func() error {
		r := o.dnsFetcher.Fetch(gCtx, u.Hostname())
		if r == nil || len(r.Records) == 0 {
			return errors.New("domain not found via DNS")
		}
		dnsResult = r
		return nil
	})

	// http fetch
	g.Go(func() error {
		select {
		case <-gCtx.Done():
			return gCtx.Err()
		default:
		}

		results := o.httpFetcher.Fetch(gCtx, *plan)
		for id, res := range results {
			httpResults[id] = res
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &domain.FetchContext{
		Target: plan.BaseURL,
		DNS:    dnsResult,
		HTTP:   httpResults,
	}, nil
}
