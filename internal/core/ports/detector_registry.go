package ports

type DetectorRegistrys struct {
	detectors []Detector
}

type DetectorRegistry interface {
	Register(detector ...Detector)
	List() []Detector
}
