package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kerjakuy/internal/dto"
	"kerjakuy/internal/middleware"
	"kerjakuy/internal/service"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	columnID, err := uuid.Parse(c.Param("columnID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid column id"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req := dto.CreateTaskRequest{ColumnID: &columnID}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ColumnID = &columnID

	task, err := h.taskService.CreateTask(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	columnID, err := uuid.Parse(c.Param("columnID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid column id"})
		return
	}

	tasks, err := h.taskService.ListTasksByColumn(c.Request.Context(), columnID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.UpdateTask(c.Request.Context(), taskID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	if err := h.taskService.DeleteTask(c.Request.Context(), taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) UpdateAssignees(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req dto.UpdateTaskAssigneesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignees, err := h.taskService.UpdateAssignees(c.Request.Context(), taskID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignees)
}

func (h *TaskHandler) AddComment(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req := dto.CreateTaskCommentRequest{TaskID: taskID}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.TaskID = taskID

	comment, err := h.taskService.AddComment(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *TaskHandler) ListComments(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	comments, err := h.taskService.ListComments(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *TaskHandler) AddAttachment(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req := dto.CreateAttachmentRequest{TaskID: taskID}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.TaskID = taskID

	attachment, err := h.taskService.AddAttachment(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func (h *TaskHandler) ListAttachments(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	attachments, err := h.taskService.ListAttachments(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attachments)
}
