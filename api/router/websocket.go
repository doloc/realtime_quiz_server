package router

import (
	"realtime_quiz_server/internal"

	"github.com/gin-gonic/gin"
)

func WebSocketRouter(apiGroup *gin.RouterGroup, hub *internal.Hub) {
	apiGroup.GET("/ws", func(c *gin.Context) {
		internal.ServeWebSocket(hub, c)
	})
}
