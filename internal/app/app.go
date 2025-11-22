package app

import (
	"kerjakuy/internal/auth"
	"kerjakuy/internal/project"
	"kerjakuy/internal/router/v1"
	"kerjakuy/internal/task"
	"kerjakuy/internal/user"
	"kerjakuy/internal/workspace"
	"kerjakuy/pkg/config"
	"kerjakuy/pkg/database"

	"github.com/gin-gonic/gin"
)

type Application struct {
	cfg *config.Config
}

func NewApplication(cfg *config.Config) *Application {
	return &Application{cfg: cfg}
}

func (a *Application) BuildRouter() *gin.Engine {
	if a.cfg.GinMode != "" {
		gin.SetMode(a.cfg.GinMode)
	}

	db := database.InitPostgresDB(a.cfg)

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	sessionRepo := auth.NewUserSessionRepository(db)
	authService := auth.NewService(userService, sessionRepo, auth.Config{
		Secret:          a.cfg.JWTSecret,
		Issuer:          a.cfg.JWTIssuer,
		AccessTokenTTL:  a.cfg.AccessTokenTTL,
		RefreshTokenTTL: a.cfg.RefreshTokenTTL,
	})

	cookieMgr := auth.NewCookieManager(auth.CookieOptions{
		AccessTTL:  a.cfg.AccessTokenTTL,
		RefreshTTL: a.cfg.RefreshTokenTTL,
	})

	authHandler := auth.NewAuthHandler(authService, cookieMgr)
	authMiddleware := auth.NewAuthMiddleware(authService)

	workspaceRepo := workspace.NewWorkspaceRepository(db)
	memberRepo := workspace.NewWorkspaceMemberRepository(db)
	workspaceService := workspace.NewWorkspaceService(workspaceRepo, memberRepo)
	workspaceHandler := workspace.NewWorkspaceHandler(workspaceService, userService)

	projectRepo := project.NewProjectRepository(db)
	boardRepo := project.NewBoardRepository(db)
	columnRepo := project.NewColumnRepository(db)
	projectService := project.NewProjectService(projectRepo, boardRepo, columnRepo)
	projectHandler := project.NewProjectHandler(projectService)

	taskRepo := task.NewTaskRepository(db)
	assigneeRepo := task.NewTaskAssigneeRepository(db)
	commentRepo := task.NewTaskCommentRepository(db)
	attachmentRepo := task.NewAttachmentRepository(db)
	taskService := task.NewService(taskRepo, assigneeRepo, commentRepo, attachmentRepo, boardRepo, columnRepo)
	taskHandler := task.NewTaskHandler(taskService)

	return router.SetupRouter(authHandler, workspaceHandler, projectHandler, taskHandler, authMiddleware)
}

func (a *Application) Run() error {
	router := a.BuildRouter()
	return router.Run(":" + a.cfg.AppPort)
}
