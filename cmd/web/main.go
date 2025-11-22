package main

import (
	"log"

	"kerjakuy/internal/app"
	"kerjakuy/pkg/config"
)

func main() {
	cfg := config.LoadConfig()
	application := app.NewApplication(cfg)

	log.Printf("Starting server on :%s\n", cfg.AppPort)

	if err := application.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
