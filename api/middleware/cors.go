package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func CORSMiddleware() gin.HandlerFunc {
// 	allowedOrigins := []string{"http://localhost:3000"}
// 	return func(c *gin.Context) {
// 		origin := c.GetHeader("Origin")
// 		for _, allowedOrigin := range allowedOrigins {
// 			if allowedOrigin == origin {
// 				c.Header("Access-Control-Allow-Origin", origin)
// 				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 				c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
// 				c.Header("Access-Control-Allow-Credentials", "true")
// 				return
// 			}
// 		}

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusNoContent)
// 			return
// 		}

// 		c.Next()
// 	}
// }

func CORSMiddleware() gin.HandlerFunc {
	// allowedOrigins := []string{"http://localhost:5173"}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// for _, allowedOrigin := range allowedOrigins {
		// 	if origin == allowedOrigin {
		// 		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		// 		break
		// 	}
		// }
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
