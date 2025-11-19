package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskComment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TaskID    uuid.UUID `gorm:"type:uuid;index" json:"task_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (tc *TaskComment) BeforeCreate(tx *gorm.DB) error {
	tc.ID = uuid.New()
	return nil
}
