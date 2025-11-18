package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Workspace struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100)" json:"name"`
	Slug      string         `gorm:"type:varchar(100);uniqueIndex" json:"slug"`
	OwnerID   uuid.UUID      `gorm:"type:uuid" json:"owner_id"`
	Plan      string         `gorm:"type:varchar(50)" json:"plan"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (w *Workspace) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New()
	return nil
}