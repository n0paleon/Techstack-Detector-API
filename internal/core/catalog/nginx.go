package catalog

import "TechstackDetectorAPI/internal/core/domain"

func Nginx(version string) *domain.Technology {
	return &domain.Technology{
		Name:        "webserver",
		Version:     version,
		Tags:        []string{"web server", "proxy", "load balancer"},
		Description: "webserver is an HTTP web server, reverse proxy, content cache, load balancer, TCP/UDP proxy server, and mail proxy server.",
		Link:        "https://nginx.org/",
	}
}
