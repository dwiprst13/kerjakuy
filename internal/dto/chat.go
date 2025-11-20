package dto

import (
	"time"

	"github.com/google/uuid"
)

type ChatChannelDTO struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	ProjectID   *uuid.UUID `json:"project_id,omitempty"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateChatChannelRequest struct {
	WorkspaceID uuid.UUID   `json:"workspace_id" binding:"required"`
	ProjectID   *uuid.UUID  `json:"project_id,omitempty"`
	Name        string      `json:"name" binding:"required,min=2,max=100"`
	Type        string      `json:"type" binding:"required,oneof=group direct"`
	UserIDs     []uuid.UUID `json:"user_ids,omitempty"`
}

type ChatChannelMemberDTO struct {
	ID        uuid.UUID `json:"id"`
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
	JoinedAt  time.Time `json:"joined_at"`
}

type AddChannelMembersRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required,dive,required"`
}

type ChatMessageDTO struct {
	ID        uuid.UUID  `json:"id"`
	ChannelID uuid.UUID  `json:"channel_id"`
	SenderID  uuid.UUID  `json:"sender_id"`
	Content   string     `json:"content"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateChatMessageRequest struct {
	ChannelID uuid.UUID  `json:"channel_id" binding:"required"`
	Content   string     `json:"content" binding:"required"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`
}

type ChatMessageReadDTO struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	ReadAt    time.Time `json:"read_at"`
}
