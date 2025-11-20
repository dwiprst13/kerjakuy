package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"kerjakuy/internal/dto"
	"kerjakuy/internal/middleware"
	authservice "kerjakuy/internal/service/auth"
)

type AuthHandler struct {
	authService authservice.Service
	cookieMgr   authservice.CookieManager
}

func NewAuthHandler(authService authservice.Service, cookieMgr authservice.CookieManager) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cookieMgr:   cookieMgr,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authService.Register(c.Request.Context(), req, h.metadataFromContext(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.handleAuthSuccess(c, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authService.Login(c.Request.Context(), req, h.metadataFromContext(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.handleAuthSuccess(c, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken, h.metadataFromContext(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.handleAuthSuccess(c, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if h.cookieMgr != nil {
		h.cookieMgr.ClearTokens(c)
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) OAuthRedirect(c *gin.Context) {
	provider := c.Param("provider")
	redirectURI := c.Query("redirect_uri")
	resp, err := h.authService.BeginOAuth(c.Request.Context(), provider, redirectURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code dan state wajib diisi"})
		return
	}
	resp, err := h.authService.HandleOAuthCallback(c.Request.Context(), provider, code, state, h.metadataFromContext(c))
	if err != nil {
		if errors.Is(err, authservice.ErrOAuthProviderNotConfigured) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) metadataFromContext(c *gin.Context) authservice.Metadata {
	return authservice.Metadata{
		UserAgent: c.Request.UserAgent(),
		IP:        c.ClientIP(),
	}
}

func (h *AuthHandler) handleAuthSuccess(c *gin.Context, status int, resp *dto.AuthResponse) {
	if h.cookieMgr != nil {
		h.cookieMgr.SetTokens(c, resp.Tokens)
	}
	c.JSON(status, resp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	userEmail, _ := middleware.GetUserEmail(c)
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   userEmail,
	})
}
