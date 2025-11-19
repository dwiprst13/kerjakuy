package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatMessageRead struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MessageID uuid.UUID `gorm:"type:uuid;index:idx_message_user,unique" json:"message_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index:idx_message_user,unique" json:"user_id"`
	ReadAt    time.Time `gorm:"autoCreateTime" json:"read_at"`
}

func (cmr *ChatMessageRead) BeforeCreate(tx *gorm.DB) error {
	cmr.ID = uuid.New()
	return nil
}
