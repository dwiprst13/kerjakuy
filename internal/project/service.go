package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"kerjakuy/internal/models"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req CreateProjectRequest, createdBy uuid.UUID) (*ProjectDTO, error)
	UpdateProject(ctx context.Context, projectID uuid.UUID, req UpdateProjectRequest) (*ProjectDTO, error)
	DeleteProject(ctx context.Context, projectID uuid.UUID) error
	ListWorkspaceProjects(ctx context.Context, workspaceID uuid.UUID) ([]ProjectDTO, error)
	CreateBoard(ctx context.Context, req CreateBoardRequest) (*BoardDTO, error)
	ListBoards(ctx context.Context, projectID uuid.UUID) ([]BoardDTO, error)
	UpdateBoard(ctx context.Context, boardID uuid.UUID, req UpdateBoardRequest) (*BoardDTO, error)
	DeleteBoard(ctx context.Context, boardID uuid.UUID) error
	CreateColumn(ctx context.Context, req CreateColumnRequest) (*ColumnDTO, error)
	ListColumns(ctx context.Context, boardID uuid.UUID) ([]ColumnDTO, error)
	UpdateColumn(ctx context.Context, columnID uuid.UUID, req UpdateColumnRequest) (*ColumnDTO, error)
	DeleteColumn(ctx context.Context, columnID uuid.UUID) error
}

type projectService struct {
	projectRepo ProjectRepository
	boardRepo   BoardRepository
	columnRepo  ColumnRepository
}

func NewProjectService(projectRepo ProjectRepository, boardRepo BoardRepository, columnRepo ColumnRepository) ProjectService {
	return &projectService{projectRepo: projectRepo, boardRepo: boardRepo, columnRepo: columnRepo}
}

func (s *projectService) CreateProject(ctx context.Context, req CreateProjectRequest, createdBy uuid.UUID) (*ProjectDTO, error) {
	project := &models.Project{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		CreatedBy:   createdBy,
	}
	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return mapProjectToDTO(project), nil
}

func (s *projectService) UpdateProject(ctx context.Context, projectID uuid.UUID, req UpdateProjectRequest) (*ProjectDTO, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = req.Description
	}
	if req.Color != nil {
		project.Color = req.Color
	}
	if req.IsArchived != nil {
		project.IsArchived = *req.IsArchived
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}
	return mapProjectToDTO(project), nil
}

func (s *projectService) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	return s.projectRepo.Delete(ctx, projectID)
}

func (s *projectService) ListWorkspaceProjects(ctx context.Context, workspaceID uuid.UUID) ([]ProjectDTO, error) {
	projects, err := s.projectRepo.ListByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	result := make([]ProjectDTO, 0, len(projects))
	for i := range projects {
		result = append(result, *mapProjectToDTO(&projects[i]))
	}
	return result, nil
}

func (s *projectService) CreateBoard(ctx context.Context, req CreateBoardRequest) (*BoardDTO, error) {
	board := &models.Board{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Position:  0,
	}
	if req.Position != nil {
		board.Position = *req.Position
	}

	if err := s.boardRepo.Create(ctx, board); err != nil {
		return nil, err
	}
	return mapBoardToDTO(board), nil
}

func (s *projectService) ListBoards(ctx context.Context, projectID uuid.UUID) ([]BoardDTO, error) {
	boards, err := s.boardRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	result := make([]BoardDTO, 0, len(boards))
	for i := range boards {
		result = append(result, *mapBoardToDTO(&boards[i]))
	}
	return result, nil
}

func (s *projectService) UpdateBoard(ctx context.Context, boardID uuid.UUID, req UpdateBoardRequest) (*BoardDTO, error) {
	board, err := s.boardRepo.FindByID(ctx, boardID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		board.Name = *req.Name
	}
	if req.Position != nil {
		board.Position = *req.Position
	}

	if err := s.boardRepo.Update(ctx, board); err != nil {
		return nil, err
	}
	return mapBoardToDTO(board), nil
}

func (s *projectService) DeleteBoard(ctx context.Context, boardID uuid.UUID) error {
	return s.boardRepo.Delete(ctx, boardID)
}

func (s *projectService) CreateColumn(ctx context.Context, req CreateColumnRequest) (*ColumnDTO, error) {
	if _, err := s.boardRepo.FindByID(ctx, req.BoardID); err != nil {
		return nil, errors.New("board not found")
	}

	column := &models.Column{
		BoardID:  req.BoardID,
		Name:     req.Name,
		Position: 0,
	}
	if req.Position != nil {
		column.Position = *req.Position
	}

	if err := s.columnRepo.Create(ctx, column); err != nil {
		return nil, err
	}
	return mapColumnToDTO(column), nil
}

func (s *projectService) ListColumns(ctx context.Context, boardID uuid.UUID) ([]ColumnDTO, error) {
	if _, err := s.boardRepo.FindByID(ctx, boardID); err != nil {
		return nil, errors.New("board not found")
	}

	columns, err := s.columnRepo.ListByBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}
	result := make([]ColumnDTO, 0, len(columns))
	for i := range columns {
		result = append(result, *mapColumnToDTO(&columns[i]))
	}
	return result, nil
}

func (s *projectService) UpdateColumn(ctx context.Context, columnID uuid.UUID, req UpdateColumnRequest) (*ColumnDTO, error) {
	column, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		column.Name = *req.Name
	}
	if req.Position != nil {
		column.Position = *req.Position
	}

	if err := s.columnRepo.Update(ctx, column); err != nil {
		return nil, err
	}
	return mapColumnToDTO(column), nil
}

func (s *projectService) DeleteColumn(ctx context.Context, columnID uuid.UUID) error {
	return s.columnRepo.Delete(ctx, columnID)
}

func mapProjectToDTO(project *models.Project) *ProjectDTO {
	return &ProjectDTO{
		ID:          project.ID,
		WorkspaceID: project.WorkspaceID,
		Name:        project.Name,
		Description: project.Description,
		Color:       project.Color,
		IsArchived:  project.IsArchived,
		CreatedBy:   project.CreatedBy,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}

func mapBoardToDTO(board *models.Board) *BoardDTO {
	return &BoardDTO{
		ID:        board.ID,
		ProjectID: board.ProjectID,
		Name:      board.Name,
		Position:  board.Position,
		CreatedAt: board.CreatedAt,
		UpdatedAt: board.UpdatedAt,
	}
}

func mapColumnToDTO(column *models.Column) *ColumnDTO {
	return &ColumnDTO{
		ID:        column.ID,
		BoardID:   column.BoardID,
		Name:      column.Name,
		Position:  column.Position,
		CreatedAt: column.CreatedAt,
		UpdatedAt: column.UpdatedAt,
	}
}
