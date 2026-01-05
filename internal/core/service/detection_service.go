package service

import (
	"context"
	"sync"

	"TechstackDetectorAPI/internal/core/domain"
	"TechstackDetectorAPI/internal/core/ports"
)

type DetectionService struct {
	registry   ports.DetectorRegistry
	fetcher    ports.Fetcher
	workerPool ports.WorkerPool
}

func NewDetectionService(
	registry ports.DetectorRegistry,
	fetcher ports.Fetcher,
	workerPool ports.WorkerPool,
) *DetectionService {
	return &DetectionService{
		registry:   registry,
		fetcher:    fetcher,
		workerPool: workerPool,
	}
}

func (s *DetectionService) Detect(
	ctx context.Context,
	target string,
) ([]domain.Technology, error) {

	detectors := s.registry.List()

	// 1️⃣ Build FetchPlan
	plan := s.buildFetchPlan(target, detectors)

	// 2️⃣ Execute fetch
	fetchCtx, err := s.fetcher.Fetch(ctx, &plan)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Run detectors concurrently
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

	// 4️⃣ Aggregate & dedup
	uniq := make(map[string]domain.Technology)

	for techs := range out {
		for _, t := range techs {
			fp := t.Fingerprint()
			uniq[fp] = t
		}
	}

	result := make([]domain.Technology, 0, len(uniq))
	for _, t := range uniq {
		result = append(result, t)
	}

	return result, nil
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

	reqMap := make(map[string]domain.FetchRequest)

	for _, d := range detectors {
		if p := d.FetchPlan(target); p != nil {
			plan.TLS = plan.TLS || p.TLS
			plan.JS = plan.JS || p.JS

			for _, r := range p.Requests {
				if _, ok := reqMap[r.ID]; !ok {
					reqMap[r.ID] = r
				}
			}
		}
	}

	for _, r := range reqMap {
		plan.Requests = append(plan.Requests, r)
	}

	return plan
}
