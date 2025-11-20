package service

import (
    "context"

    "github.com/google/uuid"
    "kerjakuy/internal/dto"
    "kerjakuy/internal/models"
    "kerjakuy/internal/repository"
)

type NotificationService interface {
    Create(ctx context.Context, userID uuid.UUID, notif dto.NotificationDTO) error
    List(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]dto.NotificationDTO, error)
    MarkRead(ctx context.Context, id uuid.UUID, isRead bool) error
}

type notificationService struct {
    repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
    return &notificationService{repo: repo}
}

func (s *notificationService) Create(ctx context.Context, userID uuid.UUID, notif dto.NotificationDTO) error {
    notification := &models.Notification{
        UserID: userID,
        Type:   notif.Type,
        Title:  notif.Title,
        Body:   notif.Body,
        Data:   notif.Data,
        IsRead: notif.IsRead,
    }
    return s.repo.Create(ctx, notification)
}

func (s *notificationService) List(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]dto.NotificationDTO, error) {
    notifications, err := s.repo.ListByUser(ctx, userID, unreadOnly)
    if err != nil {
        return nil, err
    }
    result := make([]dto.NotificationDTO, 0, len(notifications))
    for i := range notifications {
        result = append(result, dto.NotificationDTO{
            ID:        notifications[i].ID,
            UserID:    notifications[i].UserID,
            Type:      notifications[i].Type,
            Title:     notifications[i].Title,
            Body:      notifications[i].Body,
            Data:      map[string]interface{}(notifications[i].Data),
            IsRead:    notifications[i].IsRead,
            CreatedAt: notifications[i].CreatedAt,
        })
    }
    return result, nil
}

func (s *notificationService) MarkRead(ctx context.Context, id uuid.UUID, isRead bool) error {
    return s.repo.UpdateStatus(ctx, id, isRead)
}
