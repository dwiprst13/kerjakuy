package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProjectDTO struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Color       *string   `json:"color,omitempty"`
	IsArchived  bool      `json:"is_archived"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	WorkspaceID uuid.UUID `json:"workspace_id" binding:"required"`
	Name        string    `json:"name" binding:"required,min=3,max=150"`
	Description *string   `json:"description,omitempty"`
	Color       *string   `json:"color,omitempty"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=150"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	IsArchived  *bool   `json:"is_archived,omitempty"`
}

type BoardDTO struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBoardRequest struct {
	ProjectID uuid.UUID `json:"project_id" binding:"required"`
	Name      string    `json:"name" binding:"required,min=3,max=150"`
	Position  *int      `json:"position,omitempty"`
}

type UpdateBoardRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=3,max=150"`
	Position *int    `json:"position,omitempty"`
}

type ColumnDTO struct {
	ID        uuid.UUID `json:"id"`
	BoardID   uuid.UUID `json:"board_id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateColumnRequest struct {
	BoardID  uuid.UUID `json:"board_id" binding:"required"`
	Name     string    `json:"name" binding:"required,min=2,max=100"`
	Position *int      `json:"position,omitempty"`
}

type UpdateColumnRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Position *int    `json:"position,omitempty"`
}
