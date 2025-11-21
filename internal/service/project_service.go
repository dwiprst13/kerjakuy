package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"kerjakuy/internal/dto"
	"kerjakuy/internal/models"
	"kerjakuy/internal/repository"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req dto.CreateProjectRequest, createdBy uuid.UUID) (*dto.ProjectDTO, error)
	UpdateProject(ctx context.Context, projectID uuid.UUID, req dto.UpdateProjectRequest) (*dto.ProjectDTO, error)
	ListWorkspaceProjects(ctx context.Context, workspaceID uuid.UUID) ([]dto.ProjectDTO, error)
	DeleteProject(ctx context.Context, projectID uuid.UUID) error

	CreateBoard(ctx context.Context, req dto.CreateBoardRequest) (*dto.BoardDTO, error)
	UpdateBoard(ctx context.Context, boardID uuid.UUID, req dto.UpdateBoardRequest) (*dto.BoardDTO, error)
	ListBoards(ctx context.Context, projectID uuid.UUID) ([]dto.BoardDTO, error)
	DeleteBoard(ctx context.Context, boardID uuid.UUID) error

	CreateColumn(ctx context.Context, req dto.CreateColumnRequest) (*dto.ColumnDTO, error)
	UpdateColumn(ctx context.Context, columnID uuid.UUID, req dto.UpdateColumnRequest) (*dto.ColumnDTO, error)
	ListColumns(ctx context.Context, boardID uuid.UUID) ([]dto.ColumnDTO, error)
	DeleteColumn(ctx context.Context, columnID uuid.UUID) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
	boardRepo   repository.BoardRepository
	columnRepo  repository.ColumnRepository
}

func NewProjectService(projectRepo repository.ProjectRepository, boardRepo repository.BoardRepository, columnRepo repository.ColumnRepository) ProjectService {
	return &projectService{projectRepo: projectRepo, boardRepo: boardRepo, columnRepo: columnRepo}
}

func (s *projectService) CreateProject(ctx context.Context, req dto.CreateProjectRequest, createdBy uuid.UUID) (*dto.ProjectDTO, error) {
	project := &models.Project{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		CreatedBy:   createdBy,
		Description: req.Description,
		Color:       req.Color,
	}
	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	return mapProjectToDTO(project), nil
}

func (s *projectService) UpdateProject(ctx context.Context, projectID uuid.UUID, req dto.UpdateProjectRequest) (*dto.ProjectDTO, error) {
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

func (s *projectService) ListWorkspaceProjects(ctx context.Context, workspaceID uuid.UUID) ([]dto.ProjectDTO, error) {
	projects, err := s.projectRepo.ListByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.ProjectDTO, 0, len(projects))
	for i := range projects {
		res = append(res, *mapProjectToDTO(&projects[i]))
	}
	return res, nil
}

func (s *projectService) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	return s.projectRepo.Delete(ctx, projectID)
}

func (s *projectService) CreateBoard(ctx context.Context, req dto.CreateBoardRequest) (*dto.BoardDTO, error) {
	var position int
	if req.Position != nil {
		position = *req.Position
		if position < 0 {
			return nil, fmt.Errorf("position must be non-negative")
		}
	} else {
		existing, err := s.boardRepo.ListByProject(ctx, req.ProjectID)
		if err != nil {
			return nil, err
		}
		position = len(existing) + 1
	}

	board := &models.Board{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Position:  position,
	}
	if err := s.boardRepo.Create(ctx, board); err != nil {
		return nil, err
	}
	return mapBoardToDTO(board), nil
}

func (s *projectService) UpdateBoard(ctx context.Context, boardID uuid.UUID, req dto.UpdateBoardRequest) (*dto.BoardDTO, error) {
	board, err := s.boardRepo.FindByID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		board.Name = *req.Name
	}
	if req.Position != nil {
		if *req.Position < 0 {
			return nil, fmt.Errorf("position must be non-negative")
		}
		board.Position = *req.Position
	}
	if err := s.boardRepo.Update(ctx, board); err != nil {
		return nil, err
	}
	return mapBoardToDTO(board), nil
}

func (s *projectService) ListBoards(ctx context.Context, projectID uuid.UUID) ([]dto.BoardDTO, error) {
	boards, err := s.boardRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.BoardDTO, 0, len(boards))
	for i := range boards {
		res = append(res, *mapBoardToDTO(&boards[i]))
	}
	return res, nil
}

func (s *projectService) DeleteBoard(ctx context.Context, boardID uuid.UUID) error {
	return s.boardRepo.Delete(ctx, boardID)
}

func (s *projectService) CreateColumn(ctx context.Context, req dto.CreateColumnRequest) (*dto.ColumnDTO, error) {
	var position int
	if req.Position != nil {
		position = *req.Position
		if position < 0 {
			return nil, fmt.Errorf("position must be non-negative")
		}
	} else {
		existing, err := s.columnRepo.ListByBoard(ctx, req.BoardID)
		if err != nil {
			return nil, err
		}
		position = len(existing) + 1
	}

	column := &models.Column{
		BoardID:  req.BoardID,
		Name:     req.Name,
		Position: position,
	}
	if err := s.columnRepo.Create(ctx, column); err != nil {
		return nil, err
	}
	return mapColumnToDTO(column), nil
}

func (s *projectService) UpdateColumn(ctx context.Context, columnID uuid.UUID, req dto.UpdateColumnRequest) (*dto.ColumnDTO, error) {
	column, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		column.Name = *req.Name
	}
	if req.Position != nil {
		if *req.Position < 0 {
			return nil, fmt.Errorf("position must be non-negative")
		}
		column.Position = *req.Position
	}
	if err := s.columnRepo.Update(ctx, column); err != nil {
		return nil, err
	}
	return mapColumnToDTO(column), nil
}

func (s *projectService) ListColumns(ctx context.Context, boardID uuid.UUID) ([]dto.ColumnDTO, error) {
	columns, err := s.columnRepo.ListByBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.ColumnDTO, 0, len(columns))
	for i := range columns {
		res = append(res, *mapColumnToDTO(&columns[i]))
	}
	return res, nil
}

func (s *projectService) DeleteColumn(ctx context.Context, columnID uuid.UUID) error {
	return s.columnRepo.Delete(ctx, columnID)
}

func mapProjectToDTO(project *models.Project) *dto.ProjectDTO {
	return &dto.ProjectDTO{
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

func mapBoardToDTO(board *models.Board) *dto.BoardDTO {
	return &dto.BoardDTO{
		ID:        board.ID,
		ProjectID: board.ProjectID,
		Name:      board.Name,
		Position:  board.Position,
		CreatedAt: board.CreatedAt,
		UpdatedAt: board.UpdatedAt,
	}
}

func mapColumnToDTO(column *models.Column) *dto.ColumnDTO {
	return &dto.ColumnDTO{
		ID:        column.ID,
		BoardID:   column.BoardID,
		Name:      column.Name,
		Position:  column.Position,
		CreatedAt: column.CreatedAt,
		UpdatedAt: column.UpdatedAt,
	}
}
