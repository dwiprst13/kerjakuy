package workspace

import (
	"context"
	"errors"
	"strings"

	"kerjakuy/internal/auth"
	"kerjakuy/internal/models"
	"kerjakuy/internal/pkg/rbac"
	"kerjakuy/internal/repository"

	"github.com/google/uuid"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, ownerID uuid.UUID, req CreateWorkspaceRequest) (*WorkspaceDTO, error)
	UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req UpdateWorkspaceRequest) (*WorkspaceDTO, error)
	ListOwnerWorkspaces(ctx context.Context, ownerID uuid.UUID) ([]WorkspaceDTO, error)
	InviteMember(ctx context.Context, actorID uuid.UUID, workspaceID uuid.UUID, userID uuid.UUID, role string) (*WorkspaceMemberDTO, error)
	ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]WorkspaceMemberDTO, error)
	UpdateMemberRole(ctx context.Context, actorID uuid.UUID, workspaceID uuid.UUID, memberID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, actorID uuid.UUID, workspaceID uuid.UUID, userID uuid.UUID) error
}

type workspaceService struct {
	workspaceRepo     WorkspaceRepository
	memberRepo        repository.WorkspaceMemberRepository
	permissionService auth.PermissionService
}

func NewWorkspaceService(workspaceRepo WorkspaceRepository, memberRepo repository.WorkspaceMemberRepository, permissionService auth.PermissionService) WorkspaceService {
	return &workspaceService{
		workspaceRepo:     workspaceRepo,
		memberRepo:        memberRepo,
		permissionService: permissionService,
	}
}

func (s *workspaceService) CreateWorkspace(ctx context.Context, ownerID uuid.UUID, req CreateWorkspaceRequest) (*WorkspaceDTO, error) {
	slug := strings.ToLower(req.Slug)
	workspace := &models.Workspace{
		Name:    req.Name,
		Slug:    slug,
		OwnerID: ownerID,
	}
	if req.Plan != "" {
		workspace.Plan = req.Plan
	}

	if err := s.workspaceRepo.Create(ctx, workspace); err != nil {
		return nil, err
	}

	member := &models.WorkspaceMember{
		WorkspaceID: workspace.ID,
		UserID:      ownerID,
		Role:        "owner",
	}
	if err := s.memberRepo.Add(ctx, member); err != nil {
		return nil, err
	}

	return mapWorkspaceToDTO(workspace), nil
}

func (s *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req UpdateWorkspaceRequest) (*WorkspaceDTO, error) {
	// TODO: Add RBAC check for UpdateWorkspace (needs actorID passed down)
	// For now, let's focus on the member management as per plan, but ideally we should fix this too.
	// The handler doesn't pass actorID to UpdateWorkspace yet.
	// I will stick to the plan for member management first to avoid changing too many signatures at once,
	// but I should really fix UpdateWorkspace too.
	// Let's stick to the requested changes for now.

	workspace, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		workspace.Name = *req.Name
	}
	if req.Plan != nil {
		workspace.Plan = *req.Plan
	}

	if err := s.workspaceRepo.Update(ctx, workspace); err != nil {
		return nil, err
	}
	return mapWorkspaceToDTO(workspace), nil
}

func (s *workspaceService) ListOwnerWorkspaces(ctx context.Context, ownerID uuid.UUID) ([]WorkspaceDTO, error) {
	workspaces, err := s.workspaceRepo.ListByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	result := make([]WorkspaceDTO, 0, len(workspaces))
	for i := range workspaces {
		result = append(result, *mapWorkspaceToDTO(&workspaces[i]))
	}
	return result, nil
}

func (s *workspaceService) InviteMember(ctx context.Context, actorID uuid.UUID, workspaceID uuid.UUID, userID uuid.UUID, role string) (*WorkspaceMemberDTO, error) {
	allowed, err := s.permissionService.HasPermission(ctx, actorID, workspaceID, rbac.PermissionInviteMember)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errors.New("permission denied")
	}

	if role == "" {
		role = "member"
	}

	member := &models.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
	}

	if err := s.memberRepo.Add(ctx, member); err != nil {
		return nil, err
	}
	return mapWorkspaceMemberToDTO(member), nil
}

func (s *workspaceService) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]WorkspaceMemberDTO, error) {
	members, err := s.memberRepo.ListByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	result := make([]WorkspaceMemberDTO, 0, len(members))
	for i := range members {
		result = append(result, *mapWorkspaceMemberToDTO(&members[i]))
	}
	return result, nil
}

func (s *workspaceService) UpdateMemberRole(ctx context.Context, actorID uuid.UUID, workspaceID uuid.UUID, memberID uuid.UUID, role string) error {
	allowed, err := s.permissionService.HasPermission(ctx, actorID, workspaceID, rbac.PermissionUpdateMember)
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New("permission denied")
	}

	if role == "" {
		return errors.New("role is required")
	}
	return s.memberRepo.UpdateRole(ctx, memberID, role)
}

func (s *workspaceService) RemoveMember(ctx context.Context, actorID uuid.UUID, workspaceID, userID uuid.UUID) error {
	allowed, err := s.permissionService.HasPermission(ctx, actorID, workspaceID, rbac.PermissionRemoveMember)
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New("permission denied")
	}
	return s.memberRepo.Remove(ctx, workspaceID, userID)
}

func mapWorkspaceToDTO(workspace *models.Workspace) *WorkspaceDTO {
	return &WorkspaceDTO{
		ID:        workspace.ID,
		Name:      workspace.Name,
		Slug:      workspace.Slug,
		Plan:      workspace.Plan,
		OwnerID:   workspace.OwnerID,
		CreatedAt: workspace.CreatedAt,
		UpdatedAt: workspace.UpdatedAt,
	}
}

func mapWorkspaceMemberToDTO(member *models.WorkspaceMember) *WorkspaceMemberDTO {
	return &WorkspaceMemberDTO{
		ID:          member.ID,
		WorkspaceID: member.WorkspaceID,
		UserID:      member.UserID,
		Role:        member.Role,
		CreatedAt:   member.CreatedAt,
	}
}
