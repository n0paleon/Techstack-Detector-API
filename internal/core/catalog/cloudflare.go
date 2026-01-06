package catalog

import "TechstackDetectorAPI/internal/core/domain"

func CloudFlare() *domain.Technology {
	return &domain.Technology{
		Name:        "Cloudflare",
		Version:     "",
		Tags:        []string{"CDN", "proxy", "isp", "cyber security", "WAF", "DNS"},
		Description: "Cloudflare is a global cloud-based platform that provides security, performance, and reliability services for websites, apps, and networks.",
		Link:        "https://cloudflare.com/",
	}
}
