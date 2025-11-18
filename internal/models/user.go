package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(100)" json:"name"`
	Email         string         `gorm:"type:varchar(150);uniqueIndex" json:"email"`
	PasswordHash  string         `gorm:"column:password_hash" json:"-"`
	Avatar        *string        `gorm:"type:text" json:"avatar,omitempty"`

	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
