package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"kerjakuy/internal/dto"
	"kerjakuy/internal/models"
	"kerjakuy/internal/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, req dto.CreateTaskRequest, createdBy uuid.UUID) (*dto.TaskDTO, error)
	UpdateTask(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskRequest) (*dto.TaskDTO, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	ListTasksByColumn(ctx context.Context, columnID uuid.UUID) ([]dto.TaskDTO, error)
	UpdateAssignees(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskAssigneesRequest) ([]dto.TaskAssigneeDTO, error)
	AddComment(ctx context.Context, req dto.CreateTaskCommentRequest, userID uuid.UUID) (*dto.TaskCommentDTO, error)
	ListComments(ctx context.Context, taskID uuid.UUID) ([]dto.TaskCommentDTO, error)
	AddAttachment(ctx context.Context, req dto.CreateAttachmentRequest, uploadedBy uuid.UUID) (*dto.AttachmentDTO, error)
	ListAttachments(ctx context.Context, taskID uuid.UUID) ([]dto.AttachmentDTO, error)
}

type taskService struct {
	taskRepo       repository.TaskRepository
	assigneeRepo   repository.TaskAssigneeRepository
	commentRepo    repository.TaskCommentRepository
	attachmentRepo repository.AttachmentRepository
}

func NewTaskService(taskRepo repository.TaskRepository, assigneeRepo repository.TaskAssigneeRepository, commentRepo repository.TaskCommentRepository, attachmentRepo repository.AttachmentRepository) TaskService {
	return &taskService{taskRepo: taskRepo, assigneeRepo: assigneeRepo, commentRepo: commentRepo, attachmentRepo: attachmentRepo}
}

func (s *taskService) CreateTask(ctx context.Context, req dto.CreateTaskRequest, createdBy uuid.UUID) (*dto.TaskDTO, error) {
	if req.ColumnID == nil {
		return nil, fmt.Errorf("column_id is required")
	}
	columnID := *req.ColumnID

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

func (s *taskService) UpdateTask(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskRequest) (*dto.TaskDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if req.ColumnID != nil {
		task.ColumnID = req.ColumnID
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

func (s *taskService) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	return s.taskRepo.Delete(ctx, taskID)
}

func (s *taskService) ListTasksByColumn(ctx context.Context, columnID uuid.UUID) ([]dto.TaskDTO, error) {
	tasks, err := s.taskRepo.ListByColumn(ctx, columnID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TaskDTO, 0, len(tasks))
	for i := range tasks {
		result = append(result, *mapTaskToDTO(&tasks[i]))
	}
	return result, nil
}

func (s *taskService) UpdateAssignees(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskAssigneesRequest) ([]dto.TaskAssigneeDTO, error) {
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
	result := make([]dto.TaskAssigneeDTO, 0, len(updated))
	for i := range updated {
		result = append(result, dto.TaskAssigneeDTO{
			ID:     updated[i].ID,
			TaskID: updated[i].TaskID,
			UserID: updated[i].UserID,
		})
	}
	return result, nil
}

func (s *taskService) AddComment(ctx context.Context, req dto.CreateTaskCommentRequest, userID uuid.UUID) (*dto.TaskCommentDTO, error) {
	comment := &models.TaskComment{
		TaskID:  req.TaskID,
		UserID:  userID,
		Content: req.Content,
	}
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return &dto.TaskCommentDTO{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (s *taskService) ListComments(ctx context.Context, taskID uuid.UUID) ([]dto.TaskCommentDTO, error) {
	comments, err := s.commentRepo.ListByTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TaskCommentDTO, 0, len(comments))
	for _, c := range comments {
		result = append(result, dto.TaskCommentDTO{
			ID:        c.ID,
			TaskID:    c.TaskID,
			UserID:    c.UserID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		})
	}
	return result, nil
}

func (s *taskService) AddAttachment(ctx context.Context, req dto.CreateAttachmentRequest, uploadedBy uuid.UUID) (*dto.AttachmentDTO, error) {
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

func (s *taskService) ListAttachments(ctx context.Context, taskID uuid.UUID) ([]dto.AttachmentDTO, error) {
	attachments, err := s.attachmentRepo.ListByTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.AttachmentDTO, 0, len(attachments))
	for i := range attachments {
		result = append(result, *mapAttachmentToDTO(&attachments[i]))
	}
	return result, nil
}

func mapTaskToDTO(task *models.Task) *dto.TaskDTO {
	return &dto.TaskDTO{
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

func mapAttachmentToDTO(a *models.Attachment) *dto.AttachmentDTO {
	return &dto.AttachmentDTO{
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
