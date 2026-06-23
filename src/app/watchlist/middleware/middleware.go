package middleware

import (
	"log"
	"net/http"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
	"strings"
	"time"
	"watchlist/commons"

	"github.com/gin-gonic/gin"
)

// This middleware: Logs request method & path; Measures request execution time; Logs how long request took
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader(constants.Authorization)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				constants.FieldMessage: constants.ErrHeaderNotFound,
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, constants.Bearer)

		username, err := utils.ExtractUsername(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				constants.FieldMessage: err.Error(),
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
