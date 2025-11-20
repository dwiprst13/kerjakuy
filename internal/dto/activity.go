package dto

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLogDTO struct {
	ID          uuid.UUID              `json:"id"`
	WorkspaceID uuid.UUID              `json:"workspace_id"`
	ProjectID   *uuid.UUID             `json:"project_id,omitempty"`
	UserID      *uuid.UUID             `json:"user_id,omitempty"`
	Action      string                 `json:"action"`
	TargetType  string                 `json:"target_type"`
	TargetID    *uuid.UUID             `json:"target_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}
