package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attachment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TaskID     uuid.UUID `gorm:"type:uuid;index" json:"task_id"`
	UploadedBy uuid.UUID `gorm:"type:uuid;column:uploaded_by" json:"uploaded_by"`
	FileName   string    `gorm:"type:varchar(255);column:file_name" json:"file_name"`
	FileURL    string    `gorm:"type:text;column:file_url" json:"file_url"`
	FileSize   *int64    `gorm:"type:bigint;column:file_size" json:"file_size,omitempty"`
	MimeType   *string   `gorm:"type:varchar(100);column:mime_type" json:"mime_type,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (a *Attachment) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
