package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatChannelMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ChannelID uuid.UUID `gorm:"type:uuid;index:idx_channel_user,unique" json:"channel_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index:idx_channel_user,unique" json:"user_id"`
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

func (ccm *ChatChannelMember) BeforeCreate(tx *gorm.DB) error {
	ccm.ID = uuid.New()
	return nil
}
