package controller

import (
	"net/http"
	"realtime_quiz_server/cache"
	"realtime_quiz_server/common"
	"realtime_quiz_server/configuration"
	"realtime_quiz_server/session"
	"realtime_quiz_server/token"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type accountController struct {
	store *gorm.DB
}

func NewAuthController(store *gorm.DB) *accountController {
	return &accountController{store: store}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (controller *accountController) Login(cf *configuration.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		var requestBody LoginRequest
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
			return
		}

		if (requestBody.Username != "admin") || (requestBody.Password != "123456") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid username or password",
			})
			return
		}

		payload, err := token.NewPayload(requestBody.Username, time.Hour*24)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
			return
		}

		tokenString, err := token.GenerateJWT(payload, cf.TokenSymmetricKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
			return
		}

		sessionID := session.GenerateHostSessionID(requestBody.Username)

		err = cache.StoreHostSession(sessionID, requestBody.Username, time.Now().Add(time.Hour*24).Unix())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return
		}

		c.SetCookie("session_id", sessionID, 60*60, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	}
}

// func validSession() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		sessionID, err := c.Request.Cookie("session_id")
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"message": "Unauthorized",
// 			})
// 			c.Abort()
// 			return
// 		}
// 		sessionIDValue := sessionID.Value
// 		username, err := cache.GetSession(sessionIDValue)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"message": "Unauthorized",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		c.Set("username", username)
// 		c.Next()
// 	}
// }
