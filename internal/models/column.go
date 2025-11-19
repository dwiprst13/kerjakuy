package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Column struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BoardID   uuid.UUID `gorm:"type:uuid;index" json:"board_id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Position  int       `gorm:"default:0" json:"position"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (c *Column) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}
