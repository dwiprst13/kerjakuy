package main

import (
	"log"

	"github.com/gin-gonic/gin"
	router "kerjakuy/internal/router/v1"

	"kerjakuy/internal/handler"
	"kerjakuy/internal/middleware"
	"kerjakuy/internal/repository"
	"kerjakuy/internal/service"
	authservice "kerjakuy/internal/service/auth"
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
	workspaceRepo := repository.NewWorkspaceRepository(db)
	workspaceMemberRepo := repository.NewWorkspaceMemberRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	taskAssigneeRepo := repository.NewTaskAssigneeRepository(db)
	taskCommentRepo := repository.NewTaskCommentRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)

	userService := service.NewUserService(userRepo)
	authService := authservice.NewService(userService, sessionRepo, authservice.Config{
		Secret:          cfg.JWTSecret,
		Issuer:          cfg.JWTIssuer,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	})

	workspaceService := service.NewWorkspaceService(workspaceRepo, workspaceMemberRepo)
	workspaceHandler := handler.NewWorkspaceHandler(workspaceService, userService)
	projectService := service.NewProjectService(projectRepo, boardRepo, columnRepo)
	projectHandler := handler.NewProjectHandler(projectService)
	taskService := service.NewTaskService(taskRepo, taskAssigneeRepo, taskCommentRepo, attachmentRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	cookieMgr := authservice.NewCookieManager(authservice.CookieOptions{
		AccessTTL:  cfg.AccessTokenTTL,
		RefreshTTL: cfg.RefreshTokenTTL,
	})

	authHandler := handler.NewAuthHandler(authService, cookieMgr)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := router.SetupRouter(db, authHandler, workspaceHandler, projectHandler, taskHandler, authMiddleware)

	port := cfg.AppPort

	log.Printf("Starting server on :%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
