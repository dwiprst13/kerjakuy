package task

import (
	"time"

	"github.com/google/uuid"
)

type TaskDTO struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	ColumnID    *uuid.UUID `json:"column_id,omitempty"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Position    int        `json:"position"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      string     `json:"status"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	WorkspaceID uuid.UUID  `json:"workspace_id" binding:"required"`
	ProjectID   uuid.UUID  `json:"project_id" binding:"required"`
	ColumnID    *uuid.UUID `json:"column_id,omitempty"`
	Title       string     `json:"title" binding:"required,min=3,max=200"`
	Description *string    `json:"description,omitempty"`
	Position    *int       `json:"position,omitempty"`
	Priority    *string    `json:"priority,omitempty" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
	ColumnID    *uuid.UUID `json:"column_id,omitempty"`
	Title       *string    `json:"title,omitempty" binding:"omitempty,min=3,max=200"`
	Description *string    `json:"description,omitempty"`
	Position    *int       `json:"position,omitempty"`
	Priority    *string    `json:"priority,omitempty" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      *string    `json:"status,omitempty" binding:"omitempty,oneof=todo in_progress done"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type TaskAssigneeDTO struct {
	ID     uuid.UUID `json:"id"`
	TaskID uuid.UUID `json:"task_id"`
	UserID uuid.UUID `json:"user_id"`
}

type UpdateTaskAssigneesRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required,dive,required"`
}

type TaskCommentDTO struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTaskCommentRequest struct {
	TaskID  uuid.UUID `json:"task_id" binding:"required"`
	Content string    `json:"content" binding:"required,min=1"`
}

type AttachmentDTO struct {
	ID         uuid.UUID `json:"id"`
	TaskID     uuid.UUID `json:"task_id"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
	FileName   string    `json:"file_name"`
	FileURL    string    `json:"file_url"`
	FileSize   *int64    `json:"file_size,omitempty"`
	MimeType   *string   `json:"mime_type,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateAttachmentRequest struct {
	TaskID   uuid.UUID `json:"task_id" binding:"required"`
	FileName string    `json:"file_name" binding:"required"`
	FileURL  string    `json:"file_url" binding:"required,url"`
	FileSize *int64    `json:"file_size,omitempty"`
	MimeType *string   `json:"mime_type,omitempty"`
}
