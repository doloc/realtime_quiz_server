package middleware

import (
	"fmt"
	"net/http"
	"realtime_quiz_server/configuration"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "payload"
)

func AuthMiddleware(cf *configuration.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is not provided"})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization type"})
			return
		}

		accessToken := fields[1]
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			// This is a simple example, in a real application, you should use a secret key
			return []byte(cf.TokenSymmetricKey), nil
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Error parsing token: %v", err)})
			return
		}

		// Check if the token is valid
		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			return
		}

		// Get the user info from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error getting user info from token"})
			return
		}

		ctx.Set(authorizationPayloadKey, claims)
		ctx.Next()
	}
}
