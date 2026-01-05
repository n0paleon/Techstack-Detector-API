package catalog

import "TechstackDetectorAPI/internal/core/domain"

func CloudFlare() domain.Technology {
	return domain.Technology{
		Name:        "CloudFlare",
		Version:     "",
		Tags:        []string{"CDN", "proxy", "isp", "cyber security"},
		Description: "Cloudflare is a global cloud-based platform that provides security, performance, and reliability services for websites, apps, and networks.",
		Link:        "https://cloudflare.com/",
	}
}
