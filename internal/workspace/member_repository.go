package workspace

import (
	"context"

	"kerjakuy/internal/models"
	"kerjakuy/internal/repository"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type workspaceMemberRepository struct {
	db *gorm.DB
}

func NewWorkspaceMemberRepository(db *gorm.DB) repository.WorkspaceMemberRepository {
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

func (r *workspaceMemberRepository) FindByUserAndWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) (*models.WorkspaceMember, error) {
	var member models.WorkspaceMember
	if err := r.db.WithContext(ctx).Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}
