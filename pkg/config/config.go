// pkg/config/config.go

package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	GinMode         string
	DBHost          string
	DBUser          string
	DBPass          string
	DBName          string
	DBPort          string
	DBSSL           string
	AppPort         string
	JWTSecret       string
	JWTIssuer       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Tidak ada file .env ditemukan. Menggunakan variabel lingkungan sistem.")
	}

	accessTTL := parseDurationWithDefault(os.Getenv("JWT_ACCESS_TTL"), 15*time.Minute)
	refreshTTL := parseDurationWithDefault(os.Getenv("JWT_REFRESH_TTL"), 7*24*time.Hour)

	cfg := &Config{
		GinMode:         os.Getenv("GIN_MODE"),
		DBHost:          os.Getenv("DB_HOST"),
		DBUser:          os.Getenv("DB_USER"),
		DBPass:          os.Getenv("DB_PASS"),
		DBName:          os.Getenv("DB_NAME"),
		DBPort:          os.Getenv("DB_PORT"),
		DBSSL:           os.Getenv("DB_SSL"),
		AppPort:         os.Getenv("APP_PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTIssuer:       os.Getenv("JWT_ISSUER"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	if cfg.JWTIssuer == "" {
		cfg.JWTIssuer = "kerjakuy"
	}

	if cfg.DBUser == "" || cfg.DBName == "" {
		log.Fatal("Konfigurasi database DB_USER atau DB_NAME hilang di lingkungan.")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("Konfigurasi JWT_SECRET wajib diisi untuk fitur autentikasi.")
	}

	return cfg
}

func parseDurationWithDefault(value string, fallback time.Duration) time.Duration {
	if value == "" {
		return fallback
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("peringatan: gagal parsing durasi %s, menggunakan default %s\n", value, fallback)
		return fallback
	}
	return d
}
