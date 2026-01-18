package laravel

import (
	"regexp"
	"strings"

	"TechstackDetectorAPI/internal/core/catalog"
	"TechstackDetectorAPI/internal/core/domain"

	"github.com/PuerkitoBio/goquery"
)

type LaravelDetector struct{}

func NewLaravelDetector() *LaravelDetector {
	return &LaravelDetector{}
}

func (d *LaravelDetector) ID() catalog.DetectorID {
	return catalog.Laravel
}

func (d *LaravelDetector) FetchPlan(target string) *domain.FetchPlan {
	return &domain.FetchPlan{
		BaseURL: target,
		Requests: []domain.FetchRequest{
			{
				ID:          "laravel-storage",
				Path:        "/storage/",
				Method:      "GET",
				Description: "Check Laravel storage directory",
			},
			{
				ID:          "laravel-mix",
				Path:        "/mix-manifest.json",
				Method:      "GET",
				Description: "Check for Laravel Mix manifest",
			},
		},
	}
}

func (d *LaravelDetector) Detect(ctx *domain.FetchContext) ([]domain.Technology, error) {
	detectedIndicators := 0
	version := ""

	// Check each indicator type
	if d.checkCookies(ctx) {
		detectedIndicators++
	}

	if d.checkHeaders(ctx) {
		detectedIndicators++
	}

	if ver := d.checkBodyPatterns(ctx); ver != "" {
		detectedIndicators++
		version = ver
	} else if d.checkBodyPatternsWithoutVersion(ctx) {
		detectedIndicators++
	}

	if d.checkFileExposure(ctx) {
		detectedIndicators++
	}

	if d.checkCSRFToken(ctx) {
		detectedIndicators++
	}

	if d.checkErrorPage(ctx) {
		detectedIndicators++
	}

	// If we found at least 2 indicators, mark as Laravel detected
	if detectedIndicators >= 2 {
		tech1 := catalog.NewTechnology(d.ID(), "")
		if version != "" {
			tech1.Version = version
		}
		tech2 := catalog.NewTechnology(catalog.PHP, "")
		return []domain.Technology{tech1, tech2}, nil
	}

	return nil, nil
}

// New method to check for Laravel default error page
func (d *LaravelDetector) checkErrorPage(ctx *domain.FetchContext) bool {
	for _, res := range ctx.HTTP {
		if res.Error != nil {
			continue
		}

		// Skip non-HTML responses
		contentType := res.Headers.Get("Content-Type")
		if !strings.Contains(strings.ToLower(contentType), "text/html") {
			continue
		}

		bodyStr := string(res.Body)
		if d.isLaravelErrorPage(bodyStr, ctx) {
			return true
		}
	}
	return false
}

func (d *LaravelDetector) isLaravelErrorPage(body string, ctx *domain.FetchContext) bool {
	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return false
	}

	// Check for Laravel error page structure
	detected := false

	// Check body for Laravel error page classes
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		class, exists := s.Attr("class")
		if exists {
			// Laravel error pages often have these classes
			if strings.Contains(class, "antialiased") &&
				strings.Contains(class, "min-h-screen") &&
				(strings.Contains(class, "bg-gray-100") || strings.Contains(class, "dark:bg-gray-900")) {
				detected = true
			}
		}
	})

	// Check for Laravel's error page container structure
	if !detected {
		doc.Find(".max-w-xl").Each(func(i int, s *goquery.Selection) {
			// Check for Laravel's error message structure
			if s.Find(".border-r.border-gray-400").Length() > 0 &&
				s.Find(".text-lg.text-gray-500.tracking-wider").Length() > 0 {
				detected = true
			}
		})
	}

	// Check for Tailwind CSS classes that are typical in Laravel error pages
	if !detected {
		// Look for specific Tailwind utility classes used in Laravel
		tailwindClasses := []string{
			"flex items-top justify-center",
			"sm:items-center sm:pt-0",
			"relative flex",
			"items-center pt-8 sm:justify-start sm:pt-0",
			"px-4 text-lg text-gray-500 border-r border-gray-400 tracking-wider",
			"ml-4 text-lg text-gray-500 uppercase tracking-wider",
		}

		classCount := 0
		for _, classPattern := range tailwindClasses {
			if strings.Contains(body, classPattern) {
				classCount++
			}
		}
		// If we find multiple Laravel-specific Tailwind classes, it's likely a Laravel error page
		if classCount >= 3 {
			return true
		}
	}

	return detected
}

func (d *LaravelDetector) checkCookies(ctx *domain.FetchContext) bool {
	for _, res := range ctx.HTTP {
		if res.Error != nil {
			continue
		}

		for _, cookie := range res.Headers.Values("Set-Cookie") {
			cookie = strings.ToLower(cookie)
			// Laravel session cookie
			if strings.Contains(cookie, "laravel_session") {
				return true
			}
			// Laravel XSRF token cookie
			if strings.Contains(cookie, "xsrf-token") ||
				strings.Contains(cookie, "x_csrf_token") {
				return true
			}
		}
	}
	return false
}

func (d *LaravelDetector) checkHeaders(ctx *domain.FetchContext) bool {
	for _, res := range ctx.HTTP {
		if res.Error != nil {
			continue
		}

		// Check for Laravel-specific headers
		if res.Headers.Get("X-Powered-By") == "Laravel" {
			return true
		}

		// Check server header
		if server := res.Headers.Get("Server"); strings.Contains(strings.ToLower(server), "laravel") {
			return true
		}
	}
	return false
}

func (d *LaravelDetector) checkBodyPatterns(ctx *domain.FetchContext) string {
	versionPatterns := map[string]*regexp.Regexp{
		"5.x": regexp.MustCompile(`Laravel Framework (\d+\.\d+\.\d+)`),
		"6.x": regexp.MustCompile(`laravel/framework.*v(\d+\.\d+\.\d+)`),
		"7.x": regexp.MustCompile(`laravel\/.*?\bv(\d+\.\d+\.\d+)`),
	}

	for _, res := range ctx.HTTP {
		if res.Error != nil {
			continue
		}

		bodyStr := string(res.Body)

		// Try to extract version
		for _, pattern := range versionPatterns {
			if matches := pattern.FindStringSubmatch(bodyStr); len(matches) > 1 {
				return matches[1]
			}
		}
	}
	return ""
}

func (d *LaravelDetector) checkBodyPatternsWithoutVersion(ctx *domain.FetchContext) bool {
	for _, res := range ctx.HTTP {
		if res.Error != nil {
			continue
		}

		bodyStr := string(res.Body)

		// Check for common Laravel error messages
		if strings.Contains(bodyStr, "Illuminate\\") ||
			strings.Contains(bodyStr, "Laravel\\") ||
			strings.Contains(bodyStr, "laravel.log") ||
			strings.Contains(bodyStr, "vendor/laravel") {
			return true
		}

		// Check for Blade template indicators
		if strings.Contains(bodyStr, "{{") &&
			(strings.Contains(bodyStr, "@extends") ||
				strings.Contains(bodyStr, "@section") ||
				strings.Contains(bodyStr, "@yield") ||
				strings.Contains(bodyStr, "@include")) {
			return true
		}

		// Check for common Laravel route patterns
		if strings.Contains(bodyStr, "/public/index.php") ||
			strings.Contains(bodyStr, "Route::") ||
			strings.Contains(bodyStr, "artisan") {
			return true
		}
	}
	return false
}

func (d *LaravelDetector) checkFileExposure(ctx *domain.FetchContext) bool {
	for _, res := range ctx.HTTP {
		if res.Error != nil {
			continue
		}

		// Check for .env file exposure (common in misconfigured Laravel apps)
		if strings.Contains(res.URL, ".env") && res.StatusCode == 200 {
			bodyStr := string(res.Body)
			if strings.Contains(bodyStr, "APP_NAME=") &&
				strings.Contains(bodyStr, "APP_ENV=") &&
				strings.Contains(bodyStr, "DB_") {
				return true
			}
		}

		// Check for storage directory listing
		if strings.Contains(res.URL, "/storage") &&
			(res.StatusCode == 200 || res.StatusCode == 403) {
			return true
		}

		// Check for mix-manifest.json
		if strings.Contains(res.URL, "mix-manifest.json") && res.StatusCode == 200 {
			return true
		}
	}
	return false
}

func (d *LaravelDetector) checkCSRFToken(ctx *domain.FetchContext) bool {
	for _, res := range ctx.HTTP {
		if res.Error != nil {
			continue
		}

		bodyStr := string(res.Body)

		// Check for CSRF meta tag
		if strings.Contains(bodyStr, "csrf-token") &&
			strings.Contains(bodyStr, "content=") {
			return true
		}

		// Check for CSRF field in forms
		if strings.Contains(bodyStr, "@csrf") ||
			strings.Contains(bodyStr, "csrf_field()") {
			return true
		}

		// Check for _token in forms
		if strings.Contains(bodyStr, "_token") &&
			strings.Contains(bodyStr, "name=\"_token\"") {
			return true
		}
	}
	return false
}
