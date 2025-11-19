package main

import (
	"log"
	"os"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"

	"kerjakuy/internal/models"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	err = db.AutoMigrate(
		&models.ActivityLog{},
		&models.Attachment{},
		&models.Board{},
		&models.ChatChannel{},
		&models.ChatChannelMember{},
		&models.ChatMessage{},
		&models.ChatMessageRead{},
		&models.Column{},
		&models.Notification{},
		&models.Project{},
		&models.TaskAssignee{},
		&models.TaskComment{},
		&models.Task{},
		&models.User{},
		&models.UserSession{},
		&models.WorkspaceMember{},
		&models.Workspace{},
	)

	if err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("Migration completed successfully!")
}
