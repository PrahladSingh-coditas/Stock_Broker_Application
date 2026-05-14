package middleware

import (
	"authentication/commons/constants"
	"fmt"
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
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Println("Authorization header:", authHeader)
		c.Set(constants.Token, tokenString)

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		fmt.Println("##################", claims, "#################")
		tokenExpiry, _ := claims["exp"].(float64)
		c.Set(constants.Expiry, int64(tokenExpiry))

		username, _ := claims["username"].(string)
		c.Set(constants.Username, username)

		cacheKey := fmt.Sprintf("BLACKLISTE_TOKEN:%s", tokenString)

		redisClient, _, redisError := utils.GetRedisClient(false)
		if redisError != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": constants.RedisConnectionFailedError,
			})
			return
		}

		if redisClient != nil {
			exists, err := redisClient.Exists(c, cacheKey).Result()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": constants.RedisOperationFailedError,
				})
				return
			}
			if exists > 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": constants.BlacklistTokenError,
				})
				return
			}
		}

		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)

	}
}
