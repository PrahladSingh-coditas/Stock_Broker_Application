package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)
// This middleware: Logs request method & path; Measures request execution time; Logs how long request took
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next() 		// Continue to the next middleware or actual route handler.
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)

	}
}
