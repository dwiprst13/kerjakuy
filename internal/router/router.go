
package router

import (
	"kerjakuy/internal/handler" 

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.GET("/ping", handler.PingHandler)
	}

	return router
}