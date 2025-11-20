package service

import (
    "context"
    "errors"
    "strings"

    "github.com/google/uuid"
    "kerjakuy/internal/dto"
    "kerjakuy/internal/models"
    "kerjakuy/internal/repository"
)

type WorkspaceService interface {
    CreateWorkspace(ctx context.Context, ownerID uuid.UUID, req dto.CreateWorkspaceRequest) (*dto.WorkspaceDTO, error)
    UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req dto.UpdateWorkspaceRequest) (*dto.WorkspaceDTO, error)
    ListOwnerWorkspaces(ctx context.Context, ownerID uuid.UUID) ([]dto.WorkspaceDTO, error)
    InviteMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID, role string) (*dto.WorkspaceMemberDTO, error)
    ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]dto.WorkspaceMemberDTO, error)
    UpdateMemberRole(ctx context.Context, memberID uuid.UUID, role string) error
    RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error
}

type workspaceService struct {
    workspaceRepo repository.WorkspaceRepository
    memberRepo    repository.WorkspaceMemberRepository
}

func NewWorkspaceService(workspaceRepo repository.WorkspaceRepository, memberRepo repository.WorkspaceMemberRepository) WorkspaceService {
    return &workspaceService{workspaceRepo: workspaceRepo, memberRepo: memberRepo}
}

func (s *workspaceService) CreateWorkspace(ctx context.Context, ownerID uuid.UUID, req dto.CreateWorkspaceRequest) (*dto.WorkspaceDTO, error) {
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

func (s *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req dto.UpdateWorkspaceRequest) (*dto.WorkspaceDTO, error) {
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

func (s *workspaceService) ListOwnerWorkspaces(ctx context.Context, ownerID uuid.UUID) ([]dto.WorkspaceDTO, error) {
    workspaces, err := s.workspaceRepo.ListByOwner(ctx, ownerID)
    if err != nil {
        return nil, err
    }
    result := make([]dto.WorkspaceDTO, 0, len(workspaces))
    for i := range workspaces {
        result = append(result, *mapWorkspaceToDTO(&workspaces[i]))
    }
    return result, nil
}

func (s *workspaceService) InviteMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID, role string) (*dto.WorkspaceMemberDTO, error) {
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

func (s *workspaceService) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]dto.WorkspaceMemberDTO, error) {
    members, err := s.memberRepo.ListByWorkspace(ctx, workspaceID)
    if err != nil {
        return nil, err
    }
    result := make([]dto.WorkspaceMemberDTO, 0, len(members))
    for i := range members {
        result = append(result, *mapWorkspaceMemberToDTO(&members[i]))
    }
    return result, nil
}

func (s *workspaceService) UpdateMemberRole(ctx context.Context, memberID uuid.UUID, role string) error {
    if role == "" {
        return errors.New("role is required")
    }
    return s.memberRepo.UpdateRole(ctx, memberID, role)
}

func (s *workspaceService) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
    return s.memberRepo.Remove(ctx, workspaceID, userID)
}

func mapWorkspaceToDTO(workspace *models.Workspace) *dto.WorkspaceDTO {
    return &dto.WorkspaceDTO{
        ID:        workspace.ID,
        Name:      workspace.Name,
        Slug:      workspace.Slug,
        Plan:      workspace.Plan,
        OwnerID:   workspace.OwnerID,
        CreatedAt: workspace.CreatedAt,
        UpdatedAt: workspace.UpdatedAt,
    }
}

func mapWorkspaceMemberToDTO(member *models.WorkspaceMember) *dto.WorkspaceMemberDTO {
    return &dto.WorkspaceMemberDTO{
        ID:          member.ID,
        WorkspaceID: member.WorkspaceID,
        UserID:      member.UserID,
        Role:        member.Role,
        CreatedAt:   member.CreatedAt,
    }
}
