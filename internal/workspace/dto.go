package workspace

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Plan      string    `json:"plan"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
	Slug string `json:"slug" binding:"required,min=3,max=100"`
	Plan string `json:"plan" binding:"omitempty,oneof=free standard pro"`
}

type UpdateWorkspaceRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Plan *string `json:"plan,omitempty" binding:"omitempty,oneof=free standard pro"`
}

type WorkspaceMemberDTO struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	UserID      uuid.UUID `json:"user_id"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type InviteWorkspaceMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=owner admin member"`
}

type UpdateWorkspaceMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=owner admin member"`
}
