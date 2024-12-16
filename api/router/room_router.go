package router

import (
	"realtime_quiz_server/api/middleware"
	"realtime_quiz_server/configuration"
	"realtime_quiz_server/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoomRouters(rg *gin.RouterGroup, store *gorm.DB, cf *configuration.Config) {
	roomController := controller.NewRoomController(store)

	authGroup := rg.Group("/room").Use(middleware.AuthMiddleware(cf))
	{
		authGroup.POST("/create-room", roomController.CreateRoom())
	}
}