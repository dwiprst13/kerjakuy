package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID         `gorm:"type:uuid;index" json:"user_id"`
	Type      string            `gorm:"type:varchar(50)" json:"type"`
	Title     string            `gorm:"type:varchar(150)" json:"title"`
	Body      *string           `gorm:"type:text" json:"body,omitempty"`
	Data      datatypes.JSONMap `gorm:"type:jsonb" json:"data,omitempty"`
	IsRead    bool              `gorm:"default:false;column:is_read" json:"is_read"`
	CreatedAt time.Time         `gorm:"autoCreateTime" json:"created_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	n.ID = uuid.New()
	return nil
}
