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

	authGroup := rg.Group("/quiz").Use(middleware.AuthMiddleware(cf))
	{
		authGroup.POST("/create-quiz", quizController.CreateQuiz())
		authGroup.POST("/start-quiz", quizController.StartQuiz())
	}
}
