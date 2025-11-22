package router

import (
	"kerjakuy/internal/auth"
	"kerjakuy/internal/project"
	"kerjakuy/internal/task"
	"kerjakuy/internal/workspace"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *auth.AuthHandler, workspaceHandler *workspace.WorkspaceHandler, projectHandler *project.ProjectHandler, taskHandler *task.TaskHandler, authMiddleware *auth.AuthMiddleware) *gin.Engine {
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/oauth/:provider", authHandler.OAuthRedirect)
			authGroup.GET("/oauth/:provider/callback", authHandler.OAuthCallback)
			authGroup.GET("/me", authMiddleware.RequireAuth(), authHandler.Me)
		}

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
