package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	TokenHash string    `gorm:"type:text;column:token_hash" json:"-"`
	UserAgent *string   `gorm:"type:text;column:user_agent" json:"user_agent,omitempty"`
	IPAddress *string   `gorm:"type:inet;column:ip_address" json:"ip_address,omitempty"`
	ExpiresAt time.Time `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (us *UserSession) BeforeCreate(tx *gorm.DB) error {
	us.ID = uuid.New()
	return nil
}
