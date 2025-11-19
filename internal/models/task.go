package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID  `gorm:"type:uuid;index" json:"workspace_id"`
	ProjectID   uuid.UUID  `gorm:"type:uuid;index" json:"project_id"`
	ColumnID    *uuid.UUID `gorm:"type:uuid;column:column_id" json:"column_id,omitempty"`
	Title       string     `gorm:"type:varchar(200)" json:"title"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	Position    int        `gorm:"default:0" json:"position"`
	Priority    string     `gorm:"type:varchar(20);default:medium" json:"priority"`
	DueDate     *time.Time `gorm:"column:due_date" json:"due_date,omitempty"`
	Status      string     `gorm:"type:varchar(20);default:open" json:"status"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid;column:created_by" json:"created_by"`
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
