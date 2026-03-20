package middleware

import (
	"stock_broker_application/src/app/watchlist/commons/constants" 
	"errors"
	"log"
	"net/http"
	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func WatchlistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader(genericConstants.Authorization)

		if authHeader == "" {
			c.IndentedJSON(http.StatusUnauthorized, constants.ErrHeaderMissing)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, constants.Bearer)

		token, err := utils.ParseToken(tokenString)

		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(token)

		if err != nil {
			if errors.Is(err, genericConstants.ClaimMappingFailedError) {
				c.IndentedJSON(http.StatusUnauthorized, genericConstants.ErrClaimMappingFailed)
				c.Abort()
				return
			}
			c.IndentedJSON(http.StatusUnauthorized, genericConstants.ErrTokenIsInvalid)
			c.Abort()
			return
		}

		username, ok := claims[constants.Subject].(string)

		if !ok {
			c.IndentedJSON(http.StatusUnauthorized, constants.ErrUsernameNotFoundInJWT)
			c.Abort()
			return
		}

		c.Set(constants.Username, username)

		c.Next()
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)
	}
}
