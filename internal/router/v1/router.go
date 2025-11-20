
package router

import (
	"net/http"

	"kerjakuy/internal/handler"
	"kerjakuy/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	_ = db
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.GET("/ping", handler.PingHandler)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/oauth/:provider", authHandler.OAuthRedirect)
			auth.GET("/oauth/:provider/callback", authHandler.OAuthCallback)
		}

		api.GET("/auth/me", authMiddleware.RequireAuth(), func(c *gin.Context) {
			userID, _ := middleware.GetUserID(c)
			userEmail, _ := middleware.GetUserEmail(c)
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"email":   userEmail,
			})
		})
	}

	return router
}
