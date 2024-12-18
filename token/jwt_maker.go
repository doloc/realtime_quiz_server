package token

import (
	"realtime_quiz_server/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(data interface{}, secretKey string) (string, error) {
	dataMap, err := utils.ConvertStructToMap(data)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"data": dataMap,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
