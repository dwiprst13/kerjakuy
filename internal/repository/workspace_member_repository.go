package repository

import (
    "context"

    "github.com/google/uuid"
    "kerjakuy/internal/models"

    "gorm.io/gorm"
)

type WorkspaceMemberRepository interface {
    Add(ctx context.Context, member *models.WorkspaceMember) error
    UpdateRole(ctx context.Context, memberID uuid.UUID, role string) error
    ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceMember, error)
    Remove(ctx context.Context, workspaceID, userID uuid.UUID) error
}

type workspaceMemberRepository struct {
    db *gorm.DB
}

func NewWorkspaceMemberRepository(db *gorm.DB) WorkspaceMemberRepository {
    return &workspaceMemberRepository{db: db}
}

func (r *workspaceMemberRepository) Add(ctx context.Context, member *models.WorkspaceMember) error {
    return r.db.WithContext(ctx).Create(member).Error
}

func (r *workspaceMemberRepository) UpdateRole(ctx context.Context, memberID uuid.UUID, role string) error {
    return r.db.WithContext(ctx).Model(&models.WorkspaceMember{}).WithContext(ctx).Where("id = ?", memberID).Update("role", role).Error
}

func (r *workspaceMemberRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceMember, error) {
    var members []models.WorkspaceMember
    if err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&members).Error; err != nil {
        return nil, err
    }
    return members, nil
}

func (r *workspaceMemberRepository) Remove(ctx context.Context, workspaceID, userID uuid.UUID) error {
    return r.db.WithContext(ctx).Where("workspace_id = ? AND user_id = ?", workspaceID, userID).Delete(&models.WorkspaceMember{}).Error
}
