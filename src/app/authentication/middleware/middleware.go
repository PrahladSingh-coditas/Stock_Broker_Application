package middleware

import (
	"authentication/commons/constants"
	"fmt"
	"log"
	"net/http"
	"stock_broker_application/src/utils"
	"strconv"
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
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Println("Authorization header:", authHeader)
		c.Set(constants.Token, tokenString)

		tokenExpiry, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		expiry, _ := strconv.Atoi(tokenExpiry)
		exp := int64(expiry)
		c.Set(constants.Expiry, exp)
		cacheKey := fmt.Sprintf("BLACKLISTE_TOKEN:%s", tokenString)

		redisClient, redisError := utils.GetRedisClient()
		if redisError != nil {
			log.Fatal("RedisConnectionFailedError")
			return
		}

		if redisClient != nil {
			exists, err := redisClient.Exists(c, cacheKey).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed performing operation",
				})
				c.Abort()
				return
			}
			if exists > 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Token Blacklisted",
				})
				c.Abort()
				return
			}
		}

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
