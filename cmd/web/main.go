package main

import (
	"log"

	"kerjakuy/internal/router"
	"kerjakuy/pkg/config"
	"kerjakuy/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitPostgresDB(cfg)

	r := router.SetupRouter(db)

	port := cfg.AppPort

	log.Printf("🚀 Starting server on :%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
