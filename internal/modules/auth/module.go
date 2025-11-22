package auth

import (
	"kerjakuy/internal/handler"
	"kerjakuy/internal/middleware"
	"kerjakuy/internal/repository"
	"kerjakuy/internal/service"
	authservice "kerjakuy/internal/service/auth"
	"kerjakuy/pkg/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module wires authentication components together and exposes HTTP routes and middleware.
type Module struct {
	handler    *handler.AuthHandler
	middleware *middleware.AuthMiddleware
}

type Dependencies struct {
	UserService service.UserService
}

func NewModule(cfg *config.Config, db *gorm.DB, deps Dependencies) *Module {
	sessionRepo := repository.NewUserSessionRepository(db)
	authSvc := authservice.NewService(deps.UserService, sessionRepo, authservice.Config{
		Secret:          cfg.JWTSecret,
		Issuer:          cfg.JWTIssuer,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	})

	cookieMgr := authservice.NewCookieManager(authservice.CookieOptions{
		AccessTTL:  cfg.AccessTokenTTL,
		RefreshTTL: cfg.RefreshTokenTTL,
	})

	return &Module{
		handler:    handler.NewAuthHandler(authSvc, cookieMgr),
		middleware: middleware.NewAuthMiddleware(authSvc),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", m.handler.Register)
		auth.POST("/login", m.handler.Login)
		auth.POST("/refresh", m.handler.Refresh)
		auth.POST("/logout", m.handler.Logout)
		auth.GET("/oauth/:provider", m.handler.OAuthRedirect)
		auth.GET("/oauth/:provider/callback", m.handler.OAuthCallback)
		auth.GET("/me", m.middleware.RequireAuth(), m.handler.Me)
	}
}

func (m *Module) Handler() *handler.AuthHandler {
	return m.handler
}

func (m *Module) AuthMiddleware() *middleware.AuthMiddleware {
	return m.middleware
}
