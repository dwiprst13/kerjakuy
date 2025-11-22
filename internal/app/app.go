package app

import (
	authmodule "kerjakuy/internal/modules/auth"
	projectmodule "kerjakuy/internal/modules/project"
	taskmodule "kerjakuy/internal/modules/task"
	usermodule "kerjakuy/internal/modules/user"
	workspacemodule "kerjakuy/internal/modules/workspace"
	routerv1 "kerjakuy/internal/router/v1"
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

	userModule := usermodule.NewModule(db)
	authModule := authmodule.NewModule(a.cfg, db, authmodule.Dependencies{
		UserService: userModule.Service(),
	})
	workspaceModule := workspacemodule.NewModule(db, workspacemodule.Dependencies{
		UserService:    userModule.Service(),
		AuthMiddleware: authModule.AuthMiddleware(),
	})
	projectModule := projectmodule.NewModule(db, projectmodule.Dependencies{
		AuthMiddleware: authModule.AuthMiddleware(),
	})
	taskModule := taskmodule.NewModule(db, taskmodule.Dependencies{
		AuthMiddleware: authModule.AuthMiddleware(),
	})

	return routerv1.SetupRouter(routerv1.Modules{
		Auth:      authModule,
		Workspace: workspaceModule,
		Project:   projectModule,
		Task:      taskModule,
	})
}

func (a *Application) Run() error {
	router := a.BuildRouter()
	return router.Run(":" + a.cfg.AppPort)
}
