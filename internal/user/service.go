package user

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"kerjakuy/internal/models"
)

type UserService interface {
	Register(ctx context.Context, req CreateUserRequest, hashedPassword string) (*UserDTO, error)
	CreateWithPassword(ctx context.Context, req CreateUserRequest) (*UserDTO, error)
	GetByID(ctx context.Context, id uuid.UUID) (*UserDTO, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req UpdateUserProfileRequest) (*UserDTO, error)
	List(ctx context.Context) ([]UserDTO, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(ctx context.Context, req CreateUserRequest, hashedPassword string) (*UserDTO, error) {
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return mapUserToDTO(user), nil
}

func (s *userService) CreateWithPassword(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s.Register(ctx, req, string(hashed))
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return mapUserToDTO(user), nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *userService) UpdateProfile(ctx context.Context, id uuid.UUID, req UpdateUserProfileRequest) (*UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id.String())
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		if len(*req.Name) == 0 {
			return nil, errors.New("name cannot be empty")
		}
		user.Name = *req.Name
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return mapUserToDTO(user), nil
}

func (s *userService) List(ctx context.Context) ([]UserDTO, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]UserDTO, 0, len(users))
	for i := range users {
		result = append(result, *mapUserToDTO(&users[i]))
	}
	return result, nil
}

func mapUserToDTO(user *models.User) *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
