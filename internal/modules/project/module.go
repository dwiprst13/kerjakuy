package project

import (
	"kerjakuy/internal/handler"
	"kerjakuy/internal/middleware"
	"kerjakuy/internal/repository"
	"kerjakuy/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module encapsulates project, board, and column HTTP handlers with routing.
type Module struct {
	handler *handler.ProjectHandler
	auth    *middleware.AuthMiddleware
}

type Dependencies struct {
	AuthMiddleware *middleware.AuthMiddleware
}

func NewModule(db *gorm.DB, deps Dependencies) *Module {
	projectRepo := repository.NewProjectRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	projectSvc := service.NewProjectService(projectRepo, boardRepo, columnRepo)

	return &Module{
		handler: handler.NewProjectHandler(projectSvc),
		auth:    deps.AuthMiddleware,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	if m.auth == nil {
		return
	}

	workspaceProjects := api.Group("/workspaces")
	workspaceProjects.Use(m.auth.RequireAuth())
	{
		workspaceProjects.POST("/:workspaceID/projects", m.handler.CreateProject)
		workspaceProjects.GET("/:workspaceID/projects", m.handler.ListProjects)
	}

	projects := api.Group("/projects")
	projects.Use(m.auth.RequireAuth())
	{
		projects.PUT("/:projectID", m.handler.UpdateProject)
		projects.DELETE("/:projectID", m.handler.DeleteProject)
		projects.POST("/:projectID/boards", m.handler.CreateBoard)
		projects.GET("/:projectID/boards", m.handler.ListBoards)
	}

	boards := api.Group("/boards")
	boards.Use(m.auth.RequireAuth())
	{
		boards.PUT("/:boardID", m.handler.UpdateBoard)
		boards.DELETE("/:boardID", m.handler.DeleteBoard)
		boards.POST("/:boardID/columns", m.handler.CreateColumn)
		boards.GET("/:boardID/columns", m.handler.ListColumns)
	}

	columns := api.Group("/columns")
	columns.Use(m.auth.RequireAuth())
	{
		columns.PUT("/:columnID", m.handler.UpdateColumn)
		columns.DELETE("/:columnID", m.handler.DeleteColumn)
	}
}

func (m *Module) Handler() *handler.ProjectHandler {
	return m.handler
}
