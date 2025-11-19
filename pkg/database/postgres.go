
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
		log.Fatalf("Gagal menyambung ke database: %v", err)
	}

	log.Println("Koneksi database berhasil!")
	return db
}