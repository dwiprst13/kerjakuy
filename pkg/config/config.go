// pkg/config/config.go

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBPort string
	DBSSL  string
	AppPort string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Tidak ada file .env ditemukan. Menggunakan variabel lingkungan sistem.")
	}

	cfg := &Config{
		DBHost: os.Getenv("DB_HOST"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
		DBPort: os.Getenv("DB_PORT"),
		DBSSL:  os.Getenv("DB_SSL"),
		AppPort: os.Getenv("APP_PORT"), 
	}

    if cfg.AppPort == "" {
        cfg.AppPort = "8080"
    }

	if cfg.DBUser == "" || cfg.DBName == "" {
		log.Fatal("Konfigurasi database DB_USER atau DB_NAME hilang di lingkungan.")
	}

	return cfg
}