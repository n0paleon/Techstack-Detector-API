package app

import (
	"TechstackDetectorAPI/internal/adapters/detector/cloudflare"
	"TechstackDetectorAPI/internal/adapters/detector/php"
	"TechstackDetectorAPI/internal/adapters/detector/webserver"
	"TechstackDetectorAPI/internal/adapters/detector/wordpress"
	"TechstackDetectorAPI/internal/adapters/fetcher"
	"TechstackDetectorAPI/internal/adapters/registry"
	"TechstackDetectorAPI/internal/adapters/workerpool"
	"TechstackDetectorAPI/internal/core/service"
	"TechstackDetectorAPI/internal/shared/security"
	echohttp "TechstackDetectorAPI/internal/transport/http/echo"
	"TechstackDetectorAPI/internal/transport/http/handler"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

func BootstrapDetectionService() *service.DetectionService {
	// resty V2 HTTP client
	httpClient := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()). // no redirect
		SetTimeout(10 * time.Second).
		SetCookieJar(nil) // disable cookie storage (DO NOT REMOVE!!!)

	// fetchers
	httpFetcher := fetcher.NewHTTPFetcher(httpClient, 100)
	dnsFetcher := fetcher.NewDNSFetcher("1.1.1.1:53", 10)

	orchestrator := fetcher.New(
		httpFetcher,
		dnsFetcher,
	)

	// worker pool
	pool := workerpool.NewAntsPool(100)

	// detector registry
	reg := registry.NewDetectorRegistry()

	reg.Register(wordpress.NewWordPressDetector())
	reg.Register(cloudflare.NewCloudFlare())
	reg.Register(webserver.NewNginx())
	reg.Register(webserver.NewApacheHTTPD())
	reg.Register(php.NewPHPDetector())
	// TODO: implement new detector

	// security layer
	rules, err := security.LoadBlacklistFile("blacklist.txt")
	if err != nil {
		panic(err)
	}
	blacklist := security.NewBlacklist(rules)
	targetValidator := security.NewTargetValidator(blacklist, dnsFetcher)

	// detection service
	return service.NewDetectionService(
		reg,
		orchestrator,
		pool,
		targetValidator,
	)
}

func BootstrapHTTPServer(detectionService *service.DetectionService) *echo.Echo {
	// http handler
	detectHandler := handler.NewDetectHandler(detectionService)

	// http server
	server := echohttp.NewServer()

	// register routes
	echohttp.RegisterRoutes(server,
		detectHandler,
	)

	return server
}
