package controller

import (
	"realtime_quiz_server/common"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type roomController struct {
	store *gorm.DB
}

func NewRoomController(store *gorm.DB) *roomController {
	return &roomController{store: store}
}

func (controller *roomController) CreateRoom() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := gonanoid.New()
		if err != nil {
			c.JSON(500, common.ErrorResponse(err))
			return
		}

		c.JSON(200, gin.H{
			"room_id": id,
		})
	}
}
