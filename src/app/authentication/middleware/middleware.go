package middleware

import (
	"authentication/commons"
	authConstants "authentication/commons/constants"
	"fmt"
	"log"
	"net/http"
	"stock_broker_application/src/utils"
	"stock_broker_application/src/constants"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// This middleware: Logs request method & path; Measures request execution time and also Logs how long request took
func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader(constants.Authorization) //it will extract header and check if its missing
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				constants.FieldMessage: constants.ErrHeaderNotFound,
			})
			c.Abort()
			return
		}

		//extracting token from bearer and setting it
		tokenString := strings.TrimPrefix(authHeader, constants.Bearer)
		c.Set(commons.Token, tokenString)

		//redis key is gonna be token string itself
		redisKey := fmt.Sprintf(authConstants.BlacklistedToken, tokenString)
		existsInRedis, err := redisClient.Exists(c.Request.Context(), redisKey).Result()
		if err != nil {
			log.Fatalf("Redis error: %v", err)
		}

		//if token exists in redis then the token is invalidated
		if existsInRedis > 0 {
			fmt.Println("Already exists") 
			c.JSON(http.StatusUnauthorized, gin.H{
				constants.FieldMessage: constants.ErrTokenInvalidated,
			})
			c.Abort()
			return
		}

		//extracting username
		username, err := utils.ExtractUsername(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				constants.FieldMessage: err.Error(),
			})
			c.Abort()
			return
		}

		//extracting token expiry
		tokenExpiry, err := utils.ExtractExpiry(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				constants.FieldMessage: err.Error(),
			})
			c.Abort()
			return
		}

		//setting username and token expiry
		c.Set(commons.Username, username)
		c.Set(commons.Expiry, tokenExpiry)

		c.Next() // Continue to the next middleware or actual route handler.
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)
	}
}
