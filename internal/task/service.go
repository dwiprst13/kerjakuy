package task

import (
	"context"
	"fmt"

	"kerjakuy/internal/auth"
	"kerjakuy/internal/models"
	"kerjakuy/internal/pkg/rbac"
	"kerjakuy/internal/project"

	"github.com/google/uuid"
)

type Service interface {
	CreateTask(ctx context.Context, req CreateTaskRequest, createdBy uuid.UUID) (*TaskDTO, error)
	UpdateTask(ctx context.Context, actorID uuid.UUID, taskID uuid.UUID, req UpdateTaskRequest) (*TaskDTO, error)
	DeleteTask(ctx context.Context, actorID uuid.UUID, taskID uuid.UUID) error
	ListTasksByColumn(ctx context.Context, columnID uuid.UUID) ([]TaskDTO, error)
	UpdateAssignees(ctx context.Context, actorID uuid.UUID, taskID uuid.UUID, req UpdateTaskAssigneesRequest) ([]TaskAssigneeDTO, error)
	AddComment(ctx context.Context, req CreateTaskCommentRequest, userID uuid.UUID) (*TaskCommentDTO, error)
	ListComments(ctx context.Context, taskID uuid.UUID) ([]TaskCommentDTO, error)
	AddAttachment(ctx context.Context, req CreateAttachmentRequest, uploadedBy uuid.UUID) (*AttachmentDTO, error)
	ListAttachments(ctx context.Context, taskID uuid.UUID) ([]AttachmentDTO, error)
}

type taskService struct {
	taskRepo          TaskRepository
	assigneeRepo      TaskAssigneeRepository
	commentRepo       TaskCommentRepository
	attachmentRepo    AttachmentRepository
	boardRepo         project.BoardRepository
	columnRepo        project.ColumnRepository
	permissionService auth.PermissionService
}

func NewService(taskRepo TaskRepository, assigneeRepo TaskAssigneeRepository, commentRepo TaskCommentRepository, attachmentRepo AttachmentRepository, boardRepo project.BoardRepository, columnRepo project.ColumnRepository, permissionService auth.PermissionService) Service {
	return &taskService{
		taskRepo:          taskRepo,
		assigneeRepo:      assigneeRepo,
		commentRepo:       commentRepo,
		attachmentRepo:    attachmentRepo,
		boardRepo:         boardRepo,
		columnRepo:        columnRepo,
		permissionService: permissionService,
	}
}

func (s *taskService) CreateTask(ctx context.Context, req CreateTaskRequest, createdBy uuid.UUID) (*TaskDTO, error) {
	allowed, err := s.permissionService.HasPermission(ctx, createdBy, req.WorkspaceID, rbac.PermissionCreateTask)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("permission denied")
	}

	if req.ColumnID == nil {
		return nil, fmt.Errorf("column_id is required")
	}
	columnID := *req.ColumnID

	column, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, fmt.Errorf("column not found")
	}
	board, err := s.boardRepo.FindByID(ctx, column.BoardID)
	if err != nil {
		return nil, fmt.Errorf("board not found for column")
	}
	if board.ProjectID != req.ProjectID {
		return nil, fmt.Errorf("column does not belong to project")
	}

	position := 0
	if req.Position != nil {
		position = *req.Position
		if position < 0 {
			return nil, fmt.Errorf("position must be non-negative")
		}
	} else {
		existing, err := s.taskRepo.ListByColumn(ctx, columnID)
		if err != nil {
			return nil, err
		}
		position = len(existing) + 1
	}

	task := &models.Task{
		WorkspaceID: req.WorkspaceID,
		ProjectID:   req.ProjectID,
		ColumnID:    &columnID,
		Title:       req.Title,
		Description: req.Description,
		CreatedBy:   createdBy,
		Position:    position,
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}
	if task.Status == "" {
		task.Status = "todo"
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	return mapTaskToDTO(task), nil
}

func (s *taskService) UpdateTask(ctx context.Context, actorID uuid.UUID, taskID uuid.UUID, req UpdateTaskRequest) (*TaskDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	allowed, err := s.permissionService.HasPermission(ctx, actorID, task.WorkspaceID, rbac.PermissionUpdateTask)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("permission denied")
	}

	if req.ColumnID != nil {
		columnID := *req.ColumnID
		column, err := s.columnRepo.FindByID(ctx, columnID)
		if err != nil {
			return nil, fmt.Errorf("column not found")
		}
		board, err := s.boardRepo.FindByID(ctx, column.BoardID)
		if err != nil {
			return nil, fmt.Errorf("board not found for column")
		}
		if board.ProjectID != task.ProjectID {
			return nil, fmt.Errorf("column does not belong to task project")
		}
		task.ColumnID = &columnID
	}
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.Position != nil {
		if *req.Position < 0 {
			return nil, fmt.Errorf("position must be non-negative")
		}
		task.Position = *req.Position
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.CompletedAt != nil {
		task.CompletedAt = req.CompletedAt
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	return mapTaskToDTO(task), nil
}

func (s *taskService) DeleteTask(ctx context.Context, actorID uuid.UUID, taskID uuid.UUID) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	allowed, err := s.permissionService.HasPermission(ctx, actorID, task.WorkspaceID, rbac.PermissionDeleteTask)
	if err != nil {
		return err
	}
	if !allowed {
		return fmt.Errorf("permission denied")
	}

	return s.taskRepo.Delete(ctx, taskID)
}

func (s *taskService) ListTasksByColumn(ctx context.Context, columnID uuid.UUID) ([]TaskDTO, error) {
	tasks, err := s.taskRepo.ListByColumn(ctx, columnID)
	if err != nil {
		return nil, err
	}
	result := make([]TaskDTO, 0, len(tasks))
	for i := range tasks {
		result = append(result, *mapTaskToDTO(&tasks[i]))
	}
	return result, nil
}

func (s *taskService) UpdateAssignees(ctx context.Context, actorID uuid.UUID, taskID uuid.UUID, req UpdateTaskAssigneesRequest) ([]TaskAssigneeDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	allowed, err := s.permissionService.HasPermission(ctx, actorID, task.WorkspaceID, rbac.PermissionUpdateTask)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("permission denied")
	}

	assignees := make([]models.TaskAssignee, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		assignees = append(assignees, models.TaskAssignee{
			TaskID: taskID,
			UserID: userID,
		})
	}

	if err := s.assigneeRepo.ReplaceAssignees(ctx, taskID, assignees); err != nil {
		return nil, err
	}

	updated, err := s.assigneeRepo.ListByTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	result := make([]TaskAssigneeDTO, 0, len(updated))
	for i := range updated {
		result = append(result, TaskAssigneeDTO{
			ID:     updated[i].ID,
			TaskID: updated[i].TaskID,
			UserID: updated[i].UserID,
		})
	}
	return result, nil
}

func (s *taskService) AddComment(ctx context.Context, req CreateTaskCommentRequest, userID uuid.UUID) (*TaskCommentDTO, error) {
	// Check if user has access to task (via workspace)
	task, err := s.taskRepo.FindByID(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}
	// Assuming any member can comment? Or use PermissionUpdateTask?
	// Let's use PermissionUpdateTask for now, or maybe a separate PermissionCommentTask.
	// rbac.go doesn't have PermissionCommentTask. I'll use PermissionUpdateTask.
	allowed, err := s.permissionService.HasPermission(ctx, userID, task.WorkspaceID, rbac.PermissionUpdateTask)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("permission denied")
	}

	comment := &models.TaskComment{
		TaskID:  req.TaskID,
		UserID:  userID,
		Content: req.Content,
	}
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return &TaskCommentDTO{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (s *taskService) ListComments(ctx context.Context, taskID uuid.UUID) ([]TaskCommentDTO, error) {
	comments, err := s.commentRepo.ListByTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	result := make([]TaskCommentDTO, 0, len(comments))
	for _, c := range comments {
		result = append(result, TaskCommentDTO{
			ID:        c.ID,
			TaskID:    c.TaskID,
			UserID:    c.UserID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		})
	}
	return result, nil
}

func (s *taskService) AddAttachment(ctx context.Context, req CreateAttachmentRequest, uploadedBy uuid.UUID) (*AttachmentDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}

	allowed, err := s.permissionService.HasPermission(ctx, uploadedBy, task.WorkspaceID, rbac.PermissionUpdateTask)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("permission denied")
	}

	attachment := &models.Attachment{
		TaskID:     req.TaskID,
		UploadedBy: uploadedBy,
		FileName:   req.FileName,
		FileURL:    req.FileURL,
		FileSize:   req.FileSize,
		MimeType:   req.MimeType,
	}
	if err := s.attachmentRepo.Create(ctx, attachment); err != nil {
		return nil, err
	}
	return mapAttachmentToDTO(attachment), nil
}

func (s *taskService) ListAttachments(ctx context.Context, taskID uuid.UUID) ([]AttachmentDTO, error) {
	attachments, err := s.attachmentRepo.ListByTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	result := make([]AttachmentDTO, 0, len(attachments))
	for i := range attachments {
		result = append(result, *mapAttachmentToDTO(&attachments[i]))
	}
	return result, nil
}

func mapTaskToDTO(task *models.Task) *TaskDTO {
	return &TaskDTO{
		ID:          task.ID,
		WorkspaceID: task.WorkspaceID,
		ProjectID:   task.ProjectID,
		ColumnID:    task.ColumnID,
		Title:       task.Title,
		Description: task.Description,
		Position:    task.Position,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		Status:      task.Status,
		CreatedBy:   task.CreatedBy,
		CompletedAt: task.CompletedAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func mapAttachmentToDTO(a *models.Attachment) *AttachmentDTO {
	return &AttachmentDTO{
		ID:         a.ID,
		TaskID:     a.TaskID,
		UploadedBy: a.UploadedBy,
		FileName:   a.FileName,
		FileURL:    a.FileURL,
		FileSize:   a.FileSize,
		MimeType:   a.MimeType,
		CreatedAt:  a.CreatedAt,
	}
}
