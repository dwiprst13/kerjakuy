package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;index" json:"project_id"`
	Name      string    `gorm:"type:varchar(150)" json:"name"`
	Position  int       `gorm:"default:0" json:"position"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (b *Board) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}
