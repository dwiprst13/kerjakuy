package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "kerjakuy/internal/dto"
    "kerjakuy/internal/middleware"
    "kerjakuy/internal/service"
)

type WorkspaceHandler struct {
    workspaceService service.WorkspaceService
    userService      service.UserService
}

func NewWorkspaceHandler(workspaceService service.WorkspaceService, userService service.UserService) *WorkspaceHandler {
    return &WorkspaceHandler{workspaceService: workspaceService, userService: userService}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
    var req dto.CreateWorkspaceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    ownerID, ok := middleware.GetUserID(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    workspace, err := h.workspaceService.CreateWorkspace(c.Request.Context(), ownerID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, workspace)
}

func (h *WorkspaceHandler) ListWorkspaces(c *gin.Context) {
    ownerID, ok := middleware.GetUserID(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    workspaces, err := h.workspaceService.ListOwnerWorkspaces(c.Request.Context(), ownerID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
    workspaceID, err := uuid.Parse(c.Param("workspaceID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
        return
    }

    var req dto.UpdateWorkspaceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    workspace, err := h.workspaceService.UpdateWorkspace(c.Request.Context(), workspaceID, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) ListMembers(c *gin.Context) {
    workspaceID, err := uuid.Parse(c.Param("workspaceID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
        return
    }

    members, err := h.workspaceService.ListMembers(c.Request.Context(), workspaceID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, members)
}

func (h *WorkspaceHandler) InviteMember(c *gin.Context) {
    workspaceID, err := uuid.Parse(c.Param("workspaceID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
        return
    }

    var req dto.InviteWorkspaceMemberRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userService.GetByEmail(c.Request.Context(), req.Email)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    member, err := h.workspaceService.InviteMember(c.Request.Context(), workspaceID, user.ID, req.Role)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, member)
}

func (h *WorkspaceHandler) UpdateMemberRole(c *gin.Context) {
    memberID, err := uuid.Parse(c.Param("memberID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
        return
    }

    var req dto.UpdateWorkspaceMemberRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.workspaceService.UpdateMemberRole(c.Request.Context(), memberID, req.Role); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func (h *WorkspaceHandler) RemoveMember(c *gin.Context) {
    workspaceID, err := uuid.Parse(c.Param("workspaceID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
        return
    }

    userID, err := uuid.Parse(c.Param("userID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    if err := h.workspaceService.RemoveMember(c.Request.Context(), workspaceID, userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
