package service

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"TechstackDetectorAPI/internal/shared/util"
	"context"
	"sort"
	"sync"

	"TechstackDetectorAPI/internal/core/domain"
	"TechstackDetectorAPI/internal/core/ports"
)

type DetectionService struct {
	registry   ports.DetectorRegistry
	fetcher    ports.Fetcher
	workerPool ports.WorkerPool
	validator  domain.TargetValidator
}

func NewDetectionService(
	registry ports.DetectorRegistry,
	fetcher ports.Fetcher,
	workerPool ports.WorkerPool,
	validator domain.TargetValidator,
) *DetectionService {
	return &DetectionService{
		registry:   registry,
		fetcher:    fetcher,
		workerPool: workerPool,
		validator:  validator,
	}
}

func (s *DetectionService) Detect(
	ctx context.Context,
	target string,
) ([]domain.Technology, error) {
	// target validation
	resolved, err := s.validator.Validate(ctx, target)
	if err != nil {
		return nil, err
	}

	detectors := s.registry.List()

	// build FetchPlan
	plan := s.buildFetchPlan(resolved.RawURL, detectors)

	// execute fetch
	fetchCtx, err := s.fetcher.Fetch(ctx, &plan)
	if err != nil {
		return nil, err
	}

	// run detectors concurrently
	out := make(chan []domain.Technology)
	var wg sync.WaitGroup

	for _, d := range detectors {
		detector := d
		wg.Add(1)

		_ = s.workerPool.Submit(func() {
			defer wg.Done()

			res, err := detector.Detect(fetchCtx)
			if err == nil && len(res) > 0 {
				out <- res
			}
		})
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	// aggregate + dedup + prioritize
	return s.aggregateTechnologies(out), nil
}

func (s *DetectionService) aggregateTechnologies(
	in <-chan []domain.Technology,
) []domain.Technology {

	uniq := make(map[string]domain.Technology)

	for techs := range in {
		for _, t := range techs {
			fp := t.Fingerprint()

			existing, ok := uniq[fp]
			if !ok {
				uniq[fp] = t
				continue
			}

			if t.Score() > existing.Score() {
				uniq[fp] = t
			}
		}
	}

	result := make([]domain.Technology, 0, len(uniq))
	for _, t := range uniq {
		result = append(result, t)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score() > result[j].Score()
	})

	return result
}

func (s *DetectionService) buildFetchPlan(
	target string,
	detectors []ports.Detector,
) domain.FetchPlan {

	plan := domain.FetchPlan{
		BaseURL: target,
		Requests: []domain.FetchRequest{
			{
				ID:          "root",
				Path:        "/",
				Description: "default target route",
			},
		},
		TLS: true,
		// TODO: implement js fetch plan
	}

	reqMap := make(map[string]domain.FetchRequest) // key = executionKey

	for _, d := range detectors {
		if p := d.FetchPlan(target); p != nil {
			plan.TLS = plan.TLS || p.TLS
			plan.JS = plan.JS || p.JS

			for _, r := range p.Requests {
				key := util.ExecutionKey(target, r)
				if _, ok := reqMap[key]; !ok {
					reqMap[key] = r
				}
			}
		}
	}

	for _, r := range reqMap {
		plan.Requests = append(plan.Requests, r)
	}

	return plan
}

func (s *DetectionService) DetectorList() []domain.Technology {
	return catalog.List()
}
