package router

import (
	"kerjakuy/internal/handler"
	"kerjakuy/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, authHandler *handler.AuthHandler, workspaceHandler *handler.WorkspaceHandler, projectHandler *handler.ProjectHandler, taskHandler *handler.TaskHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	_ = db
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		// Ping endpoint
		api.GET("/ping", handler.PingHandler)

		// Auth endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/oauth/:provider", authHandler.OAuthRedirect)
			auth.GET("/oauth/:provider/callback", authHandler.OAuthCallback)
			auth.GET("/me", authMiddleware.RequireAuth(), authHandler.Me)
		}

		// Workspace endpoints
		workspaces := api.Group("/workspaces")
		workspaces.Use(authMiddleware.RequireAuth())
		{
			workspaces.POST("", workspaceHandler.CreateWorkspace)
			workspaces.GET("", workspaceHandler.ListWorkspaces)
			workspaces.PUT("/:workspaceID", workspaceHandler.UpdateWorkspace)

			workspaces.GET("/:workspaceID/members", workspaceHandler.ListMembers)
			workspaces.POST("/:workspaceID/members", workspaceHandler.InviteMember)
			workspaces.PATCH("/:workspaceID/members/:memberID", workspaceHandler.UpdateMemberRole)
			workspaces.DELETE("/:workspaceID/members/:userID", workspaceHandler.RemoveMember)

			workspaces.POST("/:workspaceID/projects", projectHandler.CreateProject)
			workspaces.GET("/:workspaceID/projects", projectHandler.ListProjects)
		}

		// Project endpoints
		projects := api.Group("/projects")
		projects.Use(authMiddleware.RequireAuth())
		{
			projects.PUT("/:projectID", projectHandler.UpdateProject)
			projects.DELETE("/:projectID", projectHandler.DeleteProject)
			projects.POST("/:projectID/boards", projectHandler.CreateBoard)
			projects.GET("/:projectID/boards", projectHandler.ListBoards)
		}

		boards := api.Group("/boards")
		boards.Use(authMiddleware.RequireAuth())
		{
			boards.PUT("/:boardID", projectHandler.UpdateBoard)
			boards.DELETE("/:boardID", projectHandler.DeleteBoard)
			boards.POST("/:boardID/columns", projectHandler.CreateColumn)
			boards.GET("/:boardID/columns", projectHandler.ListColumns)
		}

		columns := api.Group("/columns")
		columns.Use(authMiddleware.RequireAuth())
		{
			columns.PUT("/:columnID", projectHandler.UpdateColumn)
			columns.DELETE("/:columnID", projectHandler.DeleteColumn)
			columns.POST("/:columnID/tasks", taskHandler.CreateTask)
			columns.GET("/:columnID/tasks", taskHandler.ListTasks)
		}

		tasks := api.Group("/tasks")
		tasks.Use(authMiddleware.RequireAuth())
		{
			tasks.PUT("/:taskID", taskHandler.UpdateTask)
			tasks.DELETE("/:taskID", taskHandler.DeleteTask)
			tasks.PUT("/:taskID/assignees", taskHandler.UpdateAssignees)
			tasks.POST("/:taskID/comments", taskHandler.AddComment)
			tasks.GET("/:taskID/comments", taskHandler.ListComments)
			tasks.POST("/:taskID/attachments", taskHandler.AddAttachment)
			tasks.GET("/:taskID/attachments", taskHandler.ListAttachments)
		}
	}

	return router
}
