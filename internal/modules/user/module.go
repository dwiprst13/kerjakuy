package user

import (
	"kerjakuy/internal/repository"
	"kerjakuy/internal/service"

	"gorm.io/gorm"
)

// Module assembles user-related dependencies for reuse across other modules.
type Module struct {
	userService service.UserService
}

func NewModule(db *gorm.DB) *Module {
	userRepo := repository.NewUserRepository(db)

	return &Module{
		userService: service.NewUserService(userRepo),
	}
}

func (m *Module) Service() service.UserService {
	return m.userService
}
