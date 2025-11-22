package workspace

import (
	"kerjakuy/internal/handler"
	"kerjakuy/internal/middleware"
	"kerjakuy/internal/repository"
	"kerjakuy/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module keeps workspace HTTP surface isolated from other features.
type Module struct {
	handler *handler.WorkspaceHandler
	auth    *middleware.AuthMiddleware
}

type Dependencies struct {
	UserService    service.UserService
	AuthMiddleware *middleware.AuthMiddleware
}

func NewModule(db *gorm.DB, deps Dependencies) *Module {
	workspaceRepo := repository.NewWorkspaceRepository(db)
	memberRepo := repository.NewWorkspaceMemberRepository(db)
	workspaceSvc := service.NewWorkspaceService(workspaceRepo, memberRepo)

	return &Module{
		handler: handler.NewWorkspaceHandler(workspaceSvc, deps.UserService),
		auth:    deps.AuthMiddleware,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	if m.auth == nil {
		return
	}

	workspaces := api.Group("/workspaces")
	workspaces.Use(m.auth.RequireAuth())
	{
		workspaces.POST("", m.handler.CreateWorkspace)
		workspaces.GET("", m.handler.ListWorkspaces)
		workspaces.PUT("/:workspaceID", m.handler.UpdateWorkspace)

		workspaces.GET("/:workspaceID/members", m.handler.ListMembers)
		workspaces.POST("/:workspaceID/members", m.handler.InviteMember)
		workspaces.PATCH("/:workspaceID/members/:memberID", m.handler.UpdateMemberRole)
		workspaces.DELETE("/:workspaceID/members/:userID", m.handler.RemoveMember)
	}
}

func (m *Module) Handler() *handler.WorkspaceHandler {
	return m.handler
}
