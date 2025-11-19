package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatMessage struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ChannelID uuid.UUID  `gorm:"type:uuid;index" json:"channel_id"`
	SenderID  uuid.UUID  `gorm:"type:uuid;column:sender_id" json:"sender_id"`
	Content   string     `gorm:"type:text" json:"content"`
	ReplyToID *uuid.UUID `gorm:"type:uuid;column:reply_to_id" json:"reply_to_id,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (cm *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	cm.ID = uuid.New()
	return nil
}
