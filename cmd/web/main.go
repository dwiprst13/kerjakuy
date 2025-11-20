package main

import (
	"log"

	"github.com/gin-gonic/gin"
	router "kerjakuy/internal/router/v1"

	"kerjakuy/internal/handler"
	"kerjakuy/internal/middleware"
	"kerjakuy/internal/repository"
	"kerjakuy/internal/service"
	"kerjakuy/pkg/config"
	"kerjakuy/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitPostgresDB(cfg)

	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewUserSessionRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userService, sessionRepo, service.AuthConfig{
		Secret:          cfg.JWTSecret,
		Issuer:          cfg.JWTIssuer,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	})

	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := router.SetupRouter(db, authHandler, authMiddleware)

	port := cfg.AppPort

	log.Printf("Starting server on :%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
