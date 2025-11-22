package router

import (
	"kerjakuy/internal/handler"
	authmodule "kerjakuy/internal/modules/auth"
	projectmodule "kerjakuy/internal/modules/project"
	taskmodule "kerjakuy/internal/modules/task"
	workspacemodule "kerjakuy/internal/modules/workspace"

	"github.com/gin-gonic/gin"
)

type Modules struct {
	Auth      *authmodule.Module
	Workspace *workspacemodule.Module
	Project   *projectmodule.Module
	Task      *taskmodule.Module
}

func SetupRouter(mod Modules) *gin.Engine {
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.GET("/ping", handler.PingHandler)

		if mod.Auth != nil {
			mod.Auth.RegisterRoutes(api)
		}
		if mod.Workspace != nil {
			mod.Workspace.RegisterRoutes(api)
		}
		if mod.Project != nil {
			mod.Project.RegisterRoutes(api)
		}
		if mod.Task != nil {
			mod.Task.RegisterRoutes(api)
		}
	}

	return router
}
