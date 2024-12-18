package router

import (
	"realtime_quiz_server/api/middleware"
	"realtime_quiz_server/configuration"
	"realtime_quiz_server/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func QuizRouters(rg *gin.RouterGroup, store *gorm.DB, cf *configuration.Config) {
	quizController := controller.NewQuizController(store)

	quizGroup := rg.Group("/quiz")
	{
		quizGroup.POST("/join-quiz", quizController.JoinQuiz())
		quizGroup.Use(middleware.AuthMiddleware(cf)).POST("/create-quiz", quizController.CreateQuiz())
	}
}
