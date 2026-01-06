package nginx

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"bytes"
	"strings"

	"TechstackDetectorAPI/internal/core/domain"
	"TechstackDetectorAPI/internal/core/ports"

	"github.com/PuerkitoBio/goquery"
)

type Detector struct{}

func NewNginx() ports.Detector {
	return &Detector{}
}

func (d *Detector) Name() string {
	return "nginx"
}

func (d *Detector) FetchPlan(target string) *domain.FetchPlan {
	return nil
}

func (d *Detector) Detect(ctx *domain.FetchContext) ([]domain.Technology, error) {
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
	if !strings.Contains(serverLower, "nginx") {
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

	// <title>Welcome to nginx!</title>
	title := strings.ToLower(strings.TrimSpace(doc.Find("title").First().Text()))
	if strings.Contains(title, "welcome to nginx") {
		score++
	}

	// <h1>Welcome to nginx!</h1>
	doc.Find("h1").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		text := strings.ToLower(strings.TrimSpace(s.Text()))
		if text == "welcome to nginx!" || text == "welcome to nginx" {
			score++
			return false
		}
		return true
	})

	// a href => nginx.org / nginx.com
	doc.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if !exists {
			return true
		}

		href = strings.ToLower(href)
		if strings.Contains(href, "nginx.org") || strings.Contains(href, "nginx.com") {
			score++
			return false
		}
		return true
	})

	// "Thank you for using nginx."
	doc.Find("em").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		text := strings.ToLower(s.Text())
		if strings.Contains(text, "thank you for using nginx") {
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
