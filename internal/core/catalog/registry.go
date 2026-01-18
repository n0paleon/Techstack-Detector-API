package catalog

import (
	"TechstackDetectorAPI/internal/core/domain"
)

// registry key must be same as ports.Detector Name() method
// but the key "Name" on this map can be different
var registry = map[DetectorID]domain.Technology{
	Cloudflare: {
		Name:        "Cloudflare",
		Version:     "",
		Tags:        []string{"CDN", "proxy", "isp", "cyber security", "WAF", "DNS"},
		Description: "Cloudflare is a global cloud-based platform that provides security, performance, and reliability services for websites, apps, and networks.",
		Link:        "https://cloudflare.com/",
	},
	ApacheHTTPD: {
		Name:        "Apache HTTP Server (HTTPD)",
		Version:     "",
		Tags:        []string{"web server", "proxy", "open source"},
		Description: "Apache HTTP Server is an open-source and free web server that is written by the Apache Software Foundation (ASF).",
		Link:        "https://httpd.apache.org/",
	},
	LiteSpeed: {
		Name:        "LiteSpeed",
		Version:     "",
		Tags:        []string{"web server", "proxy", "proprietary software"},
		Description: "LiteSpeed is a high-performance, high-scalability web server with comprehensive features.",
		Link:        "https://en.wikipedia.org/wiki/LiteSpeed_Web_Server",
	},
	Nginx: {
		Name:        "nginx",
		Version:     "",
		Tags:        []string{"web server", "proxy", "load balancer"},
		Description: "nginx is an HTTP web server, reverse proxy, content cache, load balancer, TCP/UDP proxy server, and mail proxy server.",
		Link:        "https://nginx.org/",
	},
	PHP: {
		Name:        "PHP",
		Version:     "",
		Tags:        []string{"programming language", "php"},
		Description: "PHP (recursive acronym for PHP: Hypertext Preprocessor) is a widely-used open source general-purpose scripting language that is especially suited for web development and can be embedded into HTML.",
		Link:        "https://www.php.net/manual/en/introduction.php",
	},
	WordPress: {
		Name:        "WordPress",
		Tags:        []string{"cms", "php", "website"},
		Version:     "",
		Description: "WordPress is an open-source CMS.",
		Link:        "https://wordpress.org/",
	},
	Laravel: {
		Name:        "Laravel",
		Tags:        []string{"web framework", "fullstack", "website development"},
		Version:     "",
		Description: "Laravel is a free and open-source PHP-based web framework for building web applications.",
		Link:        "https://laravel.com/",
	},
}
