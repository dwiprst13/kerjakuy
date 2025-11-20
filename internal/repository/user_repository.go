package repository

import (
    "context"

    "kerjakuy/internal/models"

    "gorm.io/gorm"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByID(ctx context.Context, id string) (*models.User, error)
    FindByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    List(ctx context.Context) ([]models.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
    var user models.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) List(ctx context.Context) ([]models.User, error) {
    var users []models.User
    if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
