package repository

import (
    "context"

    "github.com/google/uuid"
    "kerjakuy/internal/models"

    "gorm.io/gorm"
)

type NotificationRepository interface {
    Create(ctx context.Context, notification *models.Notification) error
    ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]models.Notification, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, isRead bool) error
}

type notificationRepository struct {
    db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
    return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
    return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]models.Notification, error) {
    var notifications []models.Notification
    query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc")
    if unreadOnly {
        query = query.Where("is_read = ?", false)
    }
    if err := query.Find(&notifications).Error; err != nil {
        return nil, err
    }
    return notifications, nil
}

func (r *notificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, isRead bool) error {
    return r.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Update("is_read", isRead).Error
}
