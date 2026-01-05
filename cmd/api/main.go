package main

import (
	"TechstackDetectorAPI/internal/app"
	"log"
)

func main() {
	detectionService := app.BootstrapDetectionService()
	server := app.BootstrapHTTPServer(detectionService)

	log.Fatal(server.Start(":8080"))
}
