package wordpress

import (
	"TechstackDetectorAPI/internal/core/catalog"
	"bytes"
	"regexp"
	"strings"

	"TechstackDetectorAPI/internal/core/domain"

	"github.com/PuerkitoBio/goquery"
)

type WordPressDetector struct{}

func NewWordPressDetector() *WordPressDetector {
	return &WordPressDetector{}
}

func (d *WordPressDetector) Name() string {
	return "wordpress"
}

func (d *WordPressDetector) FetchPlan(target string) *domain.FetchPlan {
	return &domain.FetchPlan{
		BaseURL: target,
		Requests: []domain.FetchRequest{
			{
				ID:          "root",
				Path:        "/",
				Description: "Homepage HTML",
			},
			{
				ID:          "wp-json",
				Path:        "/wp-json/",
				Description: "WordPress REST API",
			},
			{
				ID:          "license",
				Path:        "/license.txt",
				Description: "License File",
			},
		},
	}
}

func (d *WordPressDetector) Detect(
	fc *domain.FetchContext,
) ([]domain.Technology, error) {

	score := 0

	// ===== wp-json VALIDATION (MANDATORY) =====
	wpJSON, ok := fc.HTTP["wp-json"]
	if !ok || wpJSON.Error != nil {
		return nil, nil
	}

	if wpJSON.StatusCode != 200 && wpJSON.StatusCode != 401 {
		return nil, nil
	}

	ct := strings.ToLower(wpJSON.Headers.Get("content-type"))
	if !strings.Contains(ct, "application/json") {
		return nil, nil
	}

	body := strings.ToLower(string(wpJSON.Body))

	if !(strings.Contains(body, `"routes"`) ||
		strings.Contains(body, `"namespaces"`) ||
		strings.Contains(body, "wp/v2")) {
		return nil, nil
	}

	// wp-json valid → base score
	score++

	// ===== HTML MARKERS =====
	if res, ok := fc.HTTP["root"]; ok && res.Error == nil {
		html := strings.ToLower(string(res.Body))

		if strings.Contains(html, "wp-content/") ||
			strings.Contains(html, "wp-includes/") {
			score++
		}

		if strings.Contains(html, `name="generator"`) &&
			strings.Contains(html, "wordpress") {
			score++
		}

		// Header hints
		for k, v := range res.Headers {
			key := strings.ToLower(k)
			val := strings.ToLower(strings.Join(v, " "))

			if key == "link" && strings.Contains(val, "wp-json") {
				score++
			}

			if key == "x-powered-by" && strings.Contains(val, "wordpress") {
				score++
			}
		}
	}

	// ===== 3️⃣ LICENSE CHECK =====
	if res, ok := fc.HTTP["license"]; ok && res.Error == nil {
		txt := strings.ToLower(string(res.Body))
		if strings.Contains(txt, "wordpress") {
			score++
		}
	}

	// ===== FINAL DECISION =====
	// wp-json (1) + minimal 1 strong signal
	if score >= 2 {
		version := getWordPressVersion(fc.HTTP["root"].Body)

		return []domain.Technology{*catalog.WordPress(version)}, nil
	}

	return nil, nil
}

func getWordPressVersion(html []byte) string {
	if len(html) == 0 {
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return ""
	}

	meta := doc.Find(`meta[name="generator"]`).First()
	content, exists := meta.Attr("content")
	if !exists {
		return ""
	}

	re := regexp.MustCompile(`(?i)wordpress\s+([\d\.]+)`)
	match := re.FindStringSubmatch(content)
	if len(match) != 2 {
		return ""
	}

	return strings.TrimSpace(match[1])
}
