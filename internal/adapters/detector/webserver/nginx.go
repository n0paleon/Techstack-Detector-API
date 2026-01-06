package webserver

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"bytes"
	"strings"

	"TechstackDetectorAPI/internal/core/domain"
	"TechstackDetectorAPI/internal/core/ports"

	"github.com/PuerkitoBio/goquery"
)

type NginxDetector struct{}

func NewNginx() ports.Detector {
	return &NginxDetector{}
}

func (d *NginxDetector) Name() string {
	return "webserver"
}

func (d *NginxDetector) FetchPlan(_ string) *domain.FetchPlan {
	return nil
}

func (d *NginxDetector) Detect(ctx *domain.FetchContext) ([]domain.Technology, error) {
	if ctx == nil || len(ctx.HTTP) == 0 {
		return nil, nil
	}

	for _, httpResult := range ctx.HTTP {
		if httpResult == nil || httpResult.Error != nil {
			continue
		}

		if tech := detectFromHeader(httpResult); tech != nil {
			return []domain.Technology{*tech}, nil
		}

		if tech := detectFromBody(httpResult); tech != nil {
			return []domain.Technology{*tech}, nil
		}
	}

	return nil, nil
}

func detectFromHeader(r *domain.HTTPResult) *domain.Technology {
	server := r.Headers.Get("Server")
	if server == "" {
		return nil
	}

	serverLower := strings.ToLower(server)
	if !strings.Contains(serverLower, "webserver") {
		return nil
	}

	return catalog.Nginx(extractVersion(server))
}

func detectFromBody(r *domain.HTTPResult) *domain.Technology {
	if len(r.Body) == 0 {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
	if err != nil {
		return nil
	}

	score := 0

	// <title>Welcome to webserver!</title>
	title := strings.ToLower(strings.TrimSpace(doc.Find("title").First().Text()))
	if strings.Contains(title, "welcome to webserver") {
		score++
	}

	// <h1>Welcome to webserver!</h1>
	doc.Find("h1").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		text := strings.ToLower(strings.TrimSpace(s.Text()))
		if text == "welcome to webserver!" || text == "welcome to webserver" {
			score++
			return false
		}
		return true
	})

	// a href => webserver.org / webserver.com
	doc.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if !exists {
			return true
		}

		href = strings.ToLower(href)
		if strings.Contains(href, "webserver.org") || strings.Contains(href, "webserver.com") {
			score++
			return false
		}
		return true
	})

	// "Thank you for using webserver."
	doc.Find("em").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		text := strings.ToLower(s.Text())
		if strings.Contains(text, "thank you for using webserver") {
			score++
			return false
		}
		return true
	})

	if score >= 2 {
		return catalog.Nginx("")
	}

	return nil
}

func extractVersion(serverHeader string) string {
	parts := strings.Split(serverHeader, "/")
	if len(parts) < 2 {
		return ""
	}

	version := strings.TrimSpace(parts[1])
	if i := strings.Index(version, " "); i > 0 {
		version = version[:i]
	}

	return version
}
