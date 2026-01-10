package webserver

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"TechstackDetectorAPI/internal/core/domain"
	"TechstackDetectorAPI/internal/core/ports"
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ApacheHTTPDDetector struct{}

func NewApacheHTTPD() ports.Detector {
	return &ApacheHTTPDDetector{}
}

func (d *ApacheHTTPDDetector) ID() catalog.DetectorID {
	return catalog.ApacheHTTPD
}

func (d *ApacheHTTPDDetector) FetchPlan(_ string) *domain.FetchPlan {
	return nil
}

func (d *ApacheHTTPDDetector) Detect(ctx *domain.FetchContext) ([]domain.Technology, error) {
	if ctx == nil || len(ctx.HTTP) == 0 {
		return nil, nil
	}

	for _, httpResult := range ctx.HTTP {
		if httpResult == nil || httpResult.Error != nil {
			continue
		}

		if tech := d.detectFromHeader(httpResult); tech != nil {
			return []domain.Technology{*tech}, nil
		}

		if tech := d.detectFromBody(httpResult); tech != nil {
			return []domain.Technology{*tech}, nil
		}
	}

	return nil, nil
}

func (d *ApacheHTTPDDetector) detectFromHeader(r *domain.HTTPResult) *domain.Technology {
	server := r.Headers.Get("Server")
	if server == "" {
		return nil
	}

	serverLower := strings.ToLower(server)
	if !strings.Contains(serverLower, "apache") {
		return nil
	}

	result := catalog.NewTechnology(d.ID(), d.extractVersion(server))

	return &result
}

func (d *ApacheHTTPDDetector) detectFromBody(r *domain.HTTPResult) *domain.Technology {
	if len(r.Body) == 0 {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
	if err != nil {
		return nil
	}

	score := 0

	title := strings.ToLower(strings.TrimSpace(doc.Find("title").First().Text()))
	if title == "it works!" || strings.Contains(title, "apache2 ubuntu default page") {
		score += 2
	}

	bodyText := strings.ToLower(string(r.Body))
	if strings.Contains(bodyText, "powered by apache") || strings.Contains(bodyText, "apache http server") {
		score++
	}

	doc.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		href, _ := s.Attr("href")
		if strings.Contains(strings.ToLower(href), "httpd.apache.org") {
			score++
			return false
		}
		return true
	})

	if score >= 2 {
		result := catalog.NewTechnology(d.ID(), "")
		return &result
	}

	return nil
}

func (d *ApacheHTTPDDetector) extractVersion(serverHeader string) string {
	parts := strings.Split(serverHeader, "/")
	if len(parts) < 2 {
		return ""
	}

	versionPart := parts[1]
	endIndex := strings.IndexAny(versionPart, " (")
	if endIndex > 0 {
		return versionPart[:endIndex]
	}

	return strings.TrimSpace(versionPart)
}
