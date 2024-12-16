package router

import (
	"realtime_quiz_server/configuration"
	"realtime_quiz_server/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AccountRouters(rg *gin.RouterGroup, store *gorm.DB, cf *configuration.Config) {
	authController := controller.NewAuthController(store)

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", authController.Login(cf))
	}
}
