package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationDTO struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Body      *string                `json:"body,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	IsRead    bool                   `json:"is_read"`
	CreatedAt time.Time              `json:"created_at"`
}

type UpdateNotificationStatusRequest struct {
	IsRead bool `json:"is_read" binding:"required"`
}
