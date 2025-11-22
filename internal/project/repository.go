package project

import (
	"context"

	"github.com/google/uuid"
	"kerjakuy/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BoardRepository interface {
	Create(ctx context.Context, board *models.Board) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Board, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.Board, error)
	Update(ctx context.Context, board *models.Board) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ColumnRepository interface {
	Create(ctx context.Context, column *models.Column) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Column, error)
	ListByBoard(ctx context.Context, boardID uuid.UUID) ([]models.Column, error)
	Update(ctx context.Context, column *models.Column) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type projectRepository struct {
	db *gorm.DB
}

type boardRepository struct {
	db *gorm.DB
}

type columnRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func NewColumnRepository(db *gorm.DB) ColumnRepository {
	return &columnRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", id).Error
}

func (r *boardRepository) Create(ctx context.Context, board *models.Board) error {
	return r.db.WithContext(ctx).Create(board).Error
}

func (r *boardRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Board, error) {
	var board models.Board
	if err := r.db.WithContext(ctx).First(&board, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.Board, error) {
	var boards []models.Board
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("position asc").Find(&boards).Error; err != nil {
		return nil, err
	}
	return boards, nil
}

func (r *boardRepository) Update(ctx context.Context, board *models.Board) error {
	return r.db.WithContext(ctx).Save(board).Error
}

func (r *boardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Board{}, "id = ?", id).Error
}

func (r *columnRepository) Create(ctx context.Context, column *models.Column) error {
	return r.db.WithContext(ctx).Create(column).Error
}

func (r *columnRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Column, error) {
	var column models.Column
	if err := r.db.WithContext(ctx).First(&column, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &column, nil
}

func (r *columnRepository) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]models.Column, error) {
	var columns []models.Column
	if err := r.db.WithContext(ctx).Where("board_id = ?", boardID).Order("position asc").Find(&columns).Error; err != nil {
		return nil, err
	}
	return columns, nil
}

func (r *columnRepository) Update(ctx context.Context, column *models.Column) error {
	return r.db.WithContext(ctx).Save(column).Error
}

func (r *columnRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Column{}, "id = ?", id).Error
}
