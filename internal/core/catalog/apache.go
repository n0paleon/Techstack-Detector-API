package catalog

import "TechstackDetectorAPI/internal/core/domain"

func ApacheHTTPD(version string) *domain.Technology {
	return &domain.Technology{
		Name:        "Apache HTTP Server (HTTPD)",
		Version:     version,
		Tags:        []string{"web server", "proxy", "open source"},
		Description: "Apache HTTP Server is an open-source and free web server that is written by the Apache Software Foundation (ASF).",
		Link:        "https://httpd.apache.org/",
	}
}
