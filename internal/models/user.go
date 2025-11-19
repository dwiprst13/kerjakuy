package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100)" json:"name"`
	Email        string    `gorm:"type:varchar(150);uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"type:text;column:password_hash" json:"-"`
	AvatarURL    *string   `gorm:"type:text;column:avatar_url" json:"avatar_url,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
