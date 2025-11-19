package database

import (
	"fmt"
	"log"

	"kerjakuy/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSL,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mengambil koneksi SQL: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database tidak dapat dihubungi: %v", err)
	}

	log.Println("Koneksi database berhasil!")
	return db
}
