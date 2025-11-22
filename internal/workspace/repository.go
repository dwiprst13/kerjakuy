package workspace

import (
	"context"

	"github.com/google/uuid"
	"kerjakuy/internal/models"

	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *models.Workspace) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error)
	FindBySlug(ctx context.Context, slug string) (*models.Workspace, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Workspace, error)
	Update(ctx context.Context, workspace *models.Workspace) error
}

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) Create(ctx context.Context, workspace *models.Workspace) error {
	return r.db.WithContext(ctx).Create(workspace).Error
}

func (r *workspaceRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := r.db.WithContext(ctx).First(&workspace, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepository) FindBySlug(ctx context.Context, slug string) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := r.db.WithContext(ctx).First(&workspace, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	if err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&workspaces).Error; err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (r *workspaceRepository) Update(ctx context.Context, workspace *models.Workspace) error {
	return r.db.WithContext(ctx).Save(workspace).Error
}
