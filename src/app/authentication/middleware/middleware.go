package middleware

import (
	"authentication/commons/constants"
	"log"
	"net/http"
	"stock_broker_application/src/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		authHeader := c.GetHeader("Authorization")
		log.Println("Authorization header:", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Header not Found",
			})
			c.Abort()
			return
		}
		log.Println("Authorization header:", authHeader)
		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		log.Println("Authorization header:", authHeader)
		username, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set(constants.Username, username)

		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)

	}
}
