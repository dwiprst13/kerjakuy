package user

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserProfileRequest struct {
	Name      *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
