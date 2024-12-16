package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type quizController struct {
	store *gorm.DB
}

func NewQuizController(store *gorm.DB) *quizController {
	return &quizController{store: store}
}

func (controller *quizController) CreateQuiz() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Create Room Successfully",
		})
	}
}

func (controller *quizController) StartQuiz() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Start Quiz Successfully",
		})
	}
}

func (controller *quizController) GetLeaderBoard() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Get Leader Board Successfully",
		})
	}
}
