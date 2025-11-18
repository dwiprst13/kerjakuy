package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkspaceMember struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID      `gorm:"type:uuid;index:idx_workspace_user,unique" json:"workspace_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index:idx_workspace_user,unique" json:"user_id"`
	Role        string         `gorm:"type:varchar(50)" json:"role"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (wm *WorkspaceMember) BeforeCreate(tx *gorm.DB) error {
	wm.ID = uuid.New()
	return nil
}
