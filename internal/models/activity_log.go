package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ActivityLog struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID         `gorm:"type:uuid;index" json:"workspace_id"`
	ProjectID   *uuid.UUID        `gorm:"type:uuid;index" json:"project_id,omitempty"`
	UserID      *uuid.UUID        `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Action      string            `gorm:"type:varchar(100)" json:"action"`
	TargetType  string            `gorm:"type:varchar(50);column:target_type" json:"target_type"`
	TargetID    *uuid.UUID        `gorm:"type:uuid;column:target_id" json:"target_id,omitempty"`
	Metadata    datatypes.JSONMap `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
}

func (al *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	al.ID = uuid.New()
	return nil
}
