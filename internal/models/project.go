package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;index" json:"workspace_id"`
	Name        string    `gorm:"type:varchar(150)" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Color       *string   `gorm:"type:varchar(20)" json:"color,omitempty"`
	IsArchived  bool      `gorm:"default:false" json:"is_archived"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;column:created_by" json:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}
