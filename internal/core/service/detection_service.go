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
	// Map to store best version of each technology
	// Key: technology name, Value: best technology found
	techMap := make(map[string]domain.Technology)

	for techs := range in {
		for _, tech := range techs {
			current, exists := techMap[tech.Name]

			if !exists {
				// First time seeing this technology
				techMap[tech.Name] = tech
				continue
			}

			// Compare and keep the one with higher score
			if s.compareTechnologies(tech, current) > 0 {
				techMap[tech.Name] = tech
			}
		}
	}

	// Convert map to slice
	result := make([]domain.Technology, 0, len(techMap))
	for _, tech := range techMap {
		result = append(result, tech)
	}

	// Sort by score descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score() > result[j].Score()
	})

	return result
}

// compareTechnologies returns:
// 1 if tech1 is better than tech2
// 0 if they're equal
// -1 if tech2 is better than tech1
func (s *DetectionService) compareTechnologies(tech1, tech2 domain.Technology) int {
	score1 := tech1.Score()
	score2 := tech2.Score()

	// First compare by score
	if score1 > score2 {
		return 1
	} else if score1 < score2 {
		return -1
	}

	// If scores are equal, prefer the one with version
	if tech1.Version != "" && tech2.Version == "" {
		return 1
	} else if tech1.Version == "" && tech2.Version != "" {
		return -1
	}

	// If both have versions, prefer more specific version
	if tech1.Version != "" && tech2.Version != "" {
		// This is a simple comparison - you might want to implement
		// semantic version comparison here
		if len(tech1.Version) > len(tech2.Version) {
			return 1
		} else if len(tech1.Version) < len(tech2.Version) {
			return -1
		}
	}

	// If still equal, keep the first one (return 0)
	return 0
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
				Method:      "GET",
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
