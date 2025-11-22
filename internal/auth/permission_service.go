package auth

import (
	"context"

	"kerjakuy/internal/pkg/rbac"
	"kerjakuy/internal/repository"

	"github.com/google/uuid"
)

type PermissionService interface {
	HasPermission(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, perm rbac.Permission) (bool, error)
}

type permissionService struct {
	memberRepo repository.WorkspaceMemberRepository
}

func NewPermissionService(memberRepo repository.WorkspaceMemberRepository) PermissionService {
	return &permissionService{memberRepo: memberRepo}
}

func (s *permissionService) HasPermission(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, perm rbac.Permission) (bool, error) {
	member, err := s.memberRepo.FindByUserAndWorkspace(ctx, userID, workspaceID)
	if err != nil {
		return false, nil
	}

	role := rbac.Role(member.Role)
	return rbac.HasPermission(role, perm), nil
}
