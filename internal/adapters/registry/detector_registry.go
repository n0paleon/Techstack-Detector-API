package registry

import "TechstackDetectorAPI/internal/core/ports"

type DetectorRegistry struct {
	detectors []ports.Detector
}

func (d *DetectorRegistry) Register(detector ...ports.Detector) {
	d.detectors = append(d.detectors, detector...)
}

func (d *DetectorRegistry) List() []ports.Detector {
	return d.detectors
}

func NewDetectorRegistry() *DetectorRegistry {
	return &DetectorRegistry{}
}
