package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskAssignee struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TaskID uuid.UUID `gorm:"type:uuid;index:idx_task_user,unique" json:"task_id"`
	UserID uuid.UUID `gorm:"type:uuid;index:idx_task_user,unique" json:"user_id"`
}

func (ta *TaskAssignee) BeforeCreate(tx *gorm.DB) error {
	ta.ID = uuid.New()
	return nil
}
