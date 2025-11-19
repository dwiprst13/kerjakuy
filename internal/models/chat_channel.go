package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatChannel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID  `gorm:"type:uuid;index" json:"workspace_id"`
	ProjectID   *uuid.UUID `gorm:"type:uuid;index" json:"project_id,omitempty"`
	Name        string     `gorm:"type:varchar(100)" json:"name"`
	Type        string     `gorm:"type:varchar(20);default:group" json:"type"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid;column:created_by" json:"created_by"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (cc *ChatChannel) BeforeCreate(tx *gorm.DB) error {
	cc.ID = uuid.New()
	return nil
}
