package main

import (
	"log"
	"os"

	"kerjakuy/internal/router"
	"kerjakuy/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	db := database.InitPostgresDB()

	r := router.SetupRouter(db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Starting server on :%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
