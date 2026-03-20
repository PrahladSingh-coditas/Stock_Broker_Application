package middleware

import (
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
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "authorization header missing",
			})
			return
		}

		// Authorization: Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		//parsing and validating
		username, err := utils.ValidateToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Store username in context
		c.Set("username", username)

		c.Next()
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)
	}
}
