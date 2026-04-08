package middleware

import (
	"authentication/commons"
	"log"
	"net/http"
	"stock_broker_application/src/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// This middleware: Logs request method & path; Measures request execution time; Logs how long request took
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization") //it will extract header and check if its missing
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Header not Found",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		username, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}

	
		c.Set(commons.Username, username)

		c.Next() // Continue to the next middleware or actual route handler.
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)

	}
}
