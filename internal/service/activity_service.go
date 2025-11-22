package service

import (
	"context"

	"github.com/google/uuid"
	"kerjakuy/internal/dto"
	"kerjakuy/internal/models"
	"kerjakuy/internal/project"
)

type ActivityService interface {
    Log(ctx context.Context, workspaceID uuid.UUID, payload models.ActivityLog) error
    ListWorkspaceActivity(ctx context.Context, workspaceID uuid.UUID, limit int) ([]dto.ActivityLogDTO, error)
}

type activityService struct {
	repo project.ActivityLogRepository
}

func NewActivityService(repo project.ActivityLogRepository) ActivityService {
	return &activityService{repo: repo}
}

func (s *activityService) Log(ctx context.Context, workspaceID uuid.UUID, payload models.ActivityLog) error {
    payload.WorkspaceID = workspaceID
    return s.repo.Create(ctx, &payload)
}

func (s *activityService) ListWorkspaceActivity(ctx context.Context, workspaceID uuid.UUID, limit int) ([]dto.ActivityLogDTO, error) {
    logs, err := s.repo.ListByWorkspace(ctx, workspaceID, limit)
    if err != nil {
        return nil, err
    }
    result := make([]dto.ActivityLogDTO, 0, len(logs))
    for i := range logs {
        result = append(result, dto.ActivityLogDTO{
            ID:          logs[i].ID,
            WorkspaceID: logs[i].WorkspaceID,
            ProjectID:   logs[i].ProjectID,
            UserID:      logs[i].UserID,
            Action:      logs[i].Action,
            TargetType:  logs[i].TargetType,
            TargetID:    logs[i].TargetID,
            Metadata:    map[string]interface{}(logs[i].Metadata),
            CreatedAt:   logs[i].CreatedAt,
        })
    }
    return result, nil
}
