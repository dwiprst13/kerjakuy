package repository

import (
    "context"

    "github.com/google/uuid"
    "kerjakuy/internal/models"

    "gorm.io/gorm"
)

type TaskRepository interface {
    Create(ctx context.Context, task *models.Task) error
    FindByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
    ListByColumn(ctx context.Context, columnID uuid.UUID) ([]models.Task, error)
    ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.Task, error)
    Update(ctx context.Context, task *models.Task) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type TaskAssigneeRepository interface {
    ReplaceAssignees(ctx context.Context, taskID uuid.UUID, assignees []models.TaskAssignee) error
    ListByTask(ctx context.Context, taskID uuid.UUID) ([]models.TaskAssignee, error)
}

type TaskCommentRepository interface {
    Create(ctx context.Context, comment *models.TaskComment) error
    ListByTask(ctx context.Context, taskID uuid.UUID) ([]models.TaskComment, error)
}

type AttachmentRepository interface {
    Create(ctx context.Context, attachment *models.Attachment) error
    ListByTask(ctx context.Context, taskID uuid.UUID) ([]models.Attachment, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type taskRepository struct {
    db *gorm.DB
}

type taskAssigneeRepository struct {
    db *gorm.DB
}

type taskCommentRepository struct {
    db *gorm.DB
}

type attachmentRepository struct {
    db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
    return &taskRepository{db: db}
}

func NewTaskAssigneeRepository(db *gorm.DB) TaskAssigneeRepository {
    return &taskAssigneeRepository{db: db}
}

func NewTaskCommentRepository(db *gorm.DB) TaskCommentRepository {
    return &taskCommentRepository{db: db}
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
    return &attachmentRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
    return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
    var task models.Task
    if err := r.db.WithContext(ctx).First(&task, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &task, nil
}

func (r *taskRepository) ListByColumn(ctx context.Context, columnID uuid.UUID) ([]models.Task, error) {
    var tasks []models.Task
    if err := r.db.WithContext(ctx).Where("column_id = ?", columnID).Order("position asc").Find(&tasks).Error; err != nil {
        return nil, err
    }
    return tasks, nil
}

func (r *taskRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.Task, error) {
    var tasks []models.Task
    if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&tasks).Error; err != nil {
        return nil, err
    }
    return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
    return r.db.WithContext(ctx).Save(task).Error
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.Task{}, "id = ?", id).Error
}

func (r *taskAssigneeRepository) ReplaceAssignees(ctx context.Context, taskID uuid.UUID, assignees []models.TaskAssignee) error {
    tx := r.db.WithContext(ctx).Begin()

    if err := tx.Where("task_id = ?", taskID).Delete(&models.TaskAssignee{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    if len(assignees) > 0 {
        if err := tx.Create(&assignees).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    return tx.Commit().Error
}

func (r *taskAssigneeRepository) ListByTask(ctx context.Context, taskID uuid.UUID) ([]models.TaskAssignee, error) {
    var assignees []models.TaskAssignee
    if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&assignees).Error; err != nil {
        return nil, err
    }
    return assignees, nil
}

func (r *taskCommentRepository) Create(ctx context.Context, comment *models.TaskComment) error {
    return r.db.WithContext(ctx).Create(comment).Error
}

func (r *taskCommentRepository) ListByTask(ctx context.Context, taskID uuid.UUID) ([]models.TaskComment, error) {
    var comments []models.TaskComment
    if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at asc").Find(&comments).Error; err != nil {
        return nil, err
    }
    return comments, nil
}

func (r *attachmentRepository) Create(ctx context.Context, attachment *models.Attachment) error {
    return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *attachmentRepository) ListByTask(ctx context.Context, taskID uuid.UUID) ([]models.Attachment, error) {
    var attachments []models.Attachment
    if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at asc").Find(&attachments).Error; err != nil {
        return nil, err
    }
    return attachments, nil
}

func (r *attachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.Attachment{}, "id = ?", id).Error
}
