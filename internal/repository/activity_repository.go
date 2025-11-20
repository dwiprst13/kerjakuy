package repository

import (
    "context"

    "github.com/google/uuid"
    "kerjakuy/internal/models"

    "gorm.io/gorm"
)

type ActivityLogRepository interface {
    Create(ctx context.Context, log *models.ActivityLog) error
    ListByWorkspace(ctx context.Context, workspaceID uuid.UUID, limit int) ([]models.ActivityLog, error)
}

type activityLogRepository struct {
    db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
    return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, log *models.ActivityLog) error {
    return r.db.WithContext(ctx).Create(log).Error
}

func (r *activityLogRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID, limit int) ([]models.ActivityLog, error) {
    var logs []models.ActivityLog
    query := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at desc")
    if limit > 0 {
        query = query.Limit(limit)
    }
    if err := query.Find(&logs).Error; err != nil {
        return nil, err
    }
    return logs, nil
}
