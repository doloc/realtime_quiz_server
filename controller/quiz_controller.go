package controller

import (
	"fmt"
	"net/http"
	"realtime_quiz_server/cache"
	"realtime_quiz_server/entity"
	"realtime_quiz_server/entity/reference"
	"realtime_quiz_server/session"
	"realtime_quiz_server/utils"
	"time"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type quizController struct {
	store *gorm.DB
}

func NewQuizController(store *gorm.DB) *quizController {
	return &quizController{store: store}
}

// type quizConfig struct {
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	Type        string `json:"type"`
// }

// type quizAnswerOption struct {
// 	Text      string `json:"text"`
// 	IsRequire bool   `json:"isRequire"`
// }

// type quizQuestion struct {
// 	ID             string             `json:"id"`
// 	Text           string             `json:"text"`
// 	TimeLimit      int                `json:"timeLimit"`
// 	CorrectAnswers []int              `json:"correctAnswers"`
// 	Options        []quizAnswerOption `json:"options"`
// }

// type requestCreateQuiz struct {
// 	RoomId    string         `json:"roomId"`
// 	Config    quizConfig     `json:"config"`
// 	Questions []quizQuestion `json:"questions"`
// }

type QuizPayload struct {
	// RoomID string `json:"roomId"`
	Config struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Type        string `json:"type"` // PUBLIC hoặc PRIVATE
	} `json:"config"`
	Questions []struct {
		ID      string `json:"id"`
		Text    string `json:"text"`
		Options []struct {
			Text      string `json:"text"`
			IsCorrect bool   `json:"isRequired"` // `isRequired` có thể đổi thành `isCorrect`
		} `json:"options"`
		CorrectAnswers []int `json:"correctAnswers"`
		TimeLimit      int   `json:"timeLimit"`
	} `json:"questions"`
}

func (controller *quizController) CreateQuiz() func(c *gin.Context) {
	return func(c *gin.Context) {
		var payload QuizPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
			})
			return
		}

		tx := controller.store.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return
		}

		id, err := gonanoid.New()
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create Quiz
		quiz := entity.Quiz{
			ID:          id,
			Title:       payload.Config.Title,
			Description: payload.Config.Description,
			Type:        reference.PUBLIC, // Hoặc parse từ payload.Config.Type
			Status:      reference.DRAFT,
		}

		if err := tx.Create(&quiz).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tạo questions và answers
		for _, q := range payload.Questions {
			question := entity.Question{
				QuizID:       quiz.ID,
				QuestionText: q.Text,
				TimeLimit:    int32(q.TimeLimit),
			}

			if err := tx.Create(&question).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Tạo answers
			for i, opt := range q.Options {
				answer := entity.Answer{
					QuestionID: question.ID,
					AnswerText: opt.Text,
					IsCorrect:  utils.Contains(q.CorrectAnswers, i),
				}

				if err := tx.Create(&answer).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Create Room Successfully",
			"quizId":  id,
		})
	}
}

type JoinQuizPayload struct {
	QuizID   string `json:"roomId"`
	Username string `json:"username"`
}

func (controller *quizController) JoinQuiz() func(c *gin.Context) {
	return func(c *gin.Context) {
		var payload JoinQuizPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
			})
			return
		}

		fmt.Println(payload)

		// create session
		sessionID := session.GenerateHostSessionID(payload.Username)

		err := cache.StorePlayerSession(sessionID, payload.Username, time.Now().Add(time.Hour*24).Unix())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Join Quiz Successfully",
			"sessionId": sessionID,
			"quizId":    payload.QuizID,
			"username":  payload.Username,
		})
	}
}
