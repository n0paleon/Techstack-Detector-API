package app

import (
	"TechstackDetectorAPI/internal/adapters/detector/cloudflare"
	"TechstackDetectorAPI/internal/adapters/detector/webserver"
	"TechstackDetectorAPI/internal/adapters/detector/wordpress"
	"TechstackDetectorAPI/internal/adapters/fetcher"
	"TechstackDetectorAPI/internal/adapters/registry"
	"TechstackDetectorAPI/internal/adapters/workerpool"
	"TechstackDetectorAPI/internal/core/service"
	echohttp "TechstackDetectorAPI/internal/transport/http/echo"
	"TechstackDetectorAPI/internal/transport/http/handler"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

func BootstrapDetectionService() *service.DetectionService {
	// 1️⃣ HTTP client
	httpClient := resty.New().
		SetTimeout(10 * time.Second)
	//SetRetryCount() // TODO: implement advance retry mechanism

	// 2️⃣ Fetchers
	httpFetcher := fetcher.NewHTTPFetcher(httpClient, 5)
	dnsFetcher := fetcher.NewDNSFetcher("1.1.1.1:53", 10)

	orchestrator := fetcher.New(
		httpFetcher,
		dnsFetcher,
	)

	// 3️⃣ Worker pool
	pool := workerpool.NewAntsPool(100)

	// 4️⃣ Detector registry
	reg := registry.NewDetectorRegistry()

	reg.Register(wordpress.NewWordPressDetector())
	reg.Register(cloudflare.NewCloudFlare())
	reg.Register(webserver.NewNginx())
	reg.Register(webserver.NewApacheHTTPD())
	// TODO: implement new detector

	// 5️⃣ Detection service
	return service.NewDetectionService(
		reg,
		orchestrator,
		pool,
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
