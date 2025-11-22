package task

import (
	"kerjakuy/internal/handler"
	"kerjakuy/internal/middleware"
	"kerjakuy/internal/repository"
	"kerjakuy/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module encapsulates task-related HTTP endpoints.
type Module struct {
	handler *handler.TaskHandler
	auth    *middleware.AuthMiddleware
}

type Dependencies struct {
	AuthMiddleware *middleware.AuthMiddleware
}

func NewModule(db *gorm.DB, deps Dependencies) *Module {
	taskRepo := repository.NewTaskRepository(db)
	assigneeRepo := repository.NewTaskAssigneeRepository(db)
	commentRepo := repository.NewTaskCommentRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	taskSvc := service.NewTaskService(taskRepo, assigneeRepo, commentRepo, attachmentRepo, boardRepo, columnRepo)

	return &Module{
		handler: handler.NewTaskHandler(taskSvc),
		auth:    deps.AuthMiddleware,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	if m.auth == nil {
		return
	}

	columns := api.Group("/columns")
	columns.Use(m.auth.RequireAuth())
	{
		columns.POST("/:columnID/tasks", m.handler.CreateTask)
		columns.GET("/:columnID/tasks", m.handler.ListTasks)
	}

	tasks := api.Group("/tasks")
	tasks.Use(m.auth.RequireAuth())
	{
		tasks.PUT("/:taskID", m.handler.UpdateTask)
		tasks.DELETE("/:taskID", m.handler.DeleteTask)
		tasks.PUT("/:taskID/assignees", m.handler.UpdateAssignees)
		tasks.POST("/:taskID/comments", m.handler.AddComment)
		tasks.GET("/:taskID/comments", m.handler.ListComments)
		tasks.POST("/:taskID/attachments", m.handler.AddAttachment)
		tasks.GET("/:taskID/attachments", m.handler.ListAttachments)
	}
}

func (m *Module) Handler() *handler.TaskHandler {
	return m.handler
}
