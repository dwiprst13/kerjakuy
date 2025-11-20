package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "kerjakuy/internal/models"

    "gorm.io/gorm"
)

type UserSessionRepository interface {
    Create(ctx context.Context, session *models.UserSession) error
    FindByTokenHash(ctx context.Context, hash string) (*models.UserSession, error)
    DeleteByID(ctx context.Context, id uuid.UUID) error
    DeleteExpired(ctx context.Context, now time.Time) error
}

type userSessionRepository struct {
    db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) UserSessionRepository {
    return &userSessionRepository{db: db}
}

func (r *userSessionRepository) Create(ctx context.Context, session *models.UserSession) error {
    return r.db.WithContext(ctx).Create(session).Error
}

func (r *userSessionRepository) FindByTokenHash(ctx context.Context, hash string) (*models.UserSession, error) {
    var session models.UserSession
    if err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&session).Error; err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *userSessionRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.UserSession{}, "id = ?", id).Error
}

func (r *userSessionRepository) DeleteExpired(ctx context.Context, now time.Time) error {
    return r.db.WithContext(ctx).Where("expires_at <= ?", now).Delete(&models.UserSession{}).Error
}
