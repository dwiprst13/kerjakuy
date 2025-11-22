package workspace

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"kerjakuy/internal/auth"
	"kerjakuy/internal/models"
	"kerjakuy/internal/pkg/rbac"
	"kerjakuy/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	db                *gorm.DB
	workspaceRepo     WorkspaceRepository
	memberRepo        repository.WorkspaceMemberRepository
	permissionService auth.PermissionService
	logger            *slog.Logger
}

func NewWorkspaceService(db *gorm.DB, workspaceRepo WorkspaceRepository, memberRepo repository.WorkspaceMemberRepository, permissionService auth.PermissionService, logger *slog.Logger) WorkspaceService {
	return &workspaceService{
		db:                db,
		workspaceRepo:     workspaceRepo,
		memberRepo:        memberRepo,
		permissionService: permissionService,
		logger:            logger,
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

	// Start Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// TODO: Ideally repositories should accept tx, but for now we rely on the fact that
		// we are not using the tx instance inside the repo methods which is a limitation.
		// To fix this properly, we need to update repositories to support transactions.
		// For this iteration, I will use the existing repo methods but this doesn't actually use the transaction
		// unless the repo methods are updated.
		// WAIT, I need to fix this. I should create new repo instances with the tx.

		txWorkspaceRepo := NewWorkspaceRepository(tx)
		txMemberRepo := NewWorkspaceMemberRepository(tx)

		if err := txWorkspaceRepo.Create(ctx, workspace); err != nil {
			return err
		}

		member := &models.WorkspaceMember{
			WorkspaceID: workspace.ID,
			UserID:      ownerID,
			Role:        "owner",
		}
		if err := txMemberRepo.Add(ctx, member); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		s.logger.Error("failed to create workspace", "error", err, "owner_id", ownerID)
		return nil, err
	}

	s.logger.Info("workspace created", "workspace_id", workspace.ID, "owner_id", ownerID)
	return mapWorkspaceToDTO(workspace), nil
}

func (s *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req UpdateWorkspaceRequest) (*WorkspaceDTO, error) {
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
		s.logger.Error("failed to update workspace", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	s.logger.Info("workspace updated", "workspace_id", workspaceID)
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
		s.logger.Error("failed to invite member", "error", err, "workspace_id", workspaceID, "user_id", userID)
		return nil, err
	}
	s.logger.Info("member invited", "workspace_id", workspaceID, "user_id", userID, "role", role)
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
	if err := s.memberRepo.UpdateRole(ctx, memberID, role); err != nil {
		s.logger.Error("failed to update member role", "error", err, "member_id", memberID, "role", role)
		return err
	}
	s.logger.Info("member role updated", "member_id", memberID, "role", role)
	return nil
}

func (s *workspaceService) RemoveMember(ctx context.Context, actorID uuid.UUID, workspaceID, userID uuid.UUID) error {
	allowed, err := s.permissionService.HasPermission(ctx, actorID, workspaceID, rbac.PermissionRemoveMember)
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New("permission denied")
	}
	if err := s.memberRepo.Remove(ctx, workspaceID, userID); err != nil {
		s.logger.Error("failed to remove member", "error", err, "workspace_id", workspaceID, "user_id", userID)
		return err
	}
	s.logger.Info("member removed", "workspace_id", workspaceID, "user_id", userID)
	return nil
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
