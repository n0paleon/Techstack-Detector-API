package catalog

import "TechstackDetectorAPI/internal/core/domain"

func Nginx(version string) *domain.Technology {
	return &domain.Technology{
		Name:        "nginx",
		Version:     version,
		Tags:        []string{"web server", "proxy", "load balancer"},
		Description: "nginx is an HTTP web server, reverse proxy, content cache, load balancer, TCP/UDP proxy server, and mail proxy server.",
		Link:        "https://nginx.org/",
	}
}
