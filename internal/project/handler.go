package project

import (
	"net/http"

	"kerjakuy/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectService ProjectService
}

func NewProjectHandler(projectService ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("workspaceID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}

	createdBy, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req := CreateProjectRequest{WorkspaceID: workspaceID}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.WorkspaceID = workspaceID

	project, err := h.projectService.CreateProject(c.Request.Context(), req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("workspaceID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}

	projects, err := h.projectService.ListWorkspaceProjects(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), actorID, projectID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	actorID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.projectService.DeleteProject(c.Request.Context(), actorID, projectID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) CreateBoard(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	req := CreateBoardRequest{ProjectID: projectID}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ProjectID = projectID

	actorID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	board, err := h.projectService.CreateBoard(c.Request.Context(), actorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, board)
}

func (h *ProjectHandler) ListBoards(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	boards, err := h.projectService.ListBoards(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, boards)
}

func (h *ProjectHandler) UpdateBoard(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("boardID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	var req UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	board, err := h.projectService.UpdateBoard(c.Request.Context(), actorID, boardID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}

func (h *ProjectHandler) DeleteBoard(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("boardID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	actorID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.projectService.DeleteBoard(c.Request.Context(), actorID, boardID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) CreateColumn(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("boardID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	req := CreateColumnRequest{BoardID: boardID}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.BoardID = boardID

	actorID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	column, err := h.projectService.CreateColumn(c.Request.Context(), actorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, column)
}

func (h *ProjectHandler) ListColumns(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("boardID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	columns, err := h.projectService.ListColumns(c.Request.Context(), boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, columns)
}

func (h *ProjectHandler) UpdateColumn(c *gin.Context) {
	columnID, err := uuid.Parse(c.Param("columnID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid column id"})
		return
	}

	var req UpdateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	column, err := h.projectService.UpdateColumn(c.Request.Context(), actorID, columnID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, column)
}

func (h *ProjectHandler) DeleteColumn(c *gin.Context) {
	columnID, err := uuid.Parse(c.Param("columnID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid column id"})
		return
	}

	actorID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.projectService.DeleteColumn(c.Request.Context(), actorID, columnID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
