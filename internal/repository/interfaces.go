package repository

import (
	"context"

	"kerjakuy/internal/models"

	"github.com/google/uuid"
)

type WorkspaceMemberRepository interface {
	Add(ctx context.Context, member *models.WorkspaceMember) error
	UpdateRole(ctx context.Context, memberID uuid.UUID, role string) error
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceMember, error)
	Remove(ctx context.Context, workspaceID, userID uuid.UUID) error
	FindByUserAndWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) (*models.WorkspaceMember, error)
}
