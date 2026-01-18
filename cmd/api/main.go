package main

import (
	"TechstackDetectorAPI/internal/app"
	"log"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(false)
}

func main() {
	detectionService := app.BootstrapDetectionService()
	server := app.BootstrapHTTPServer(detectionService)

	log.Fatal(server.Start(":8080"))
}
