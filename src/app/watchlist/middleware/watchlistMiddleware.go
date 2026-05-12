package middleware

import (
	"stock_broker_application/src/app/watchlist/commons/constants"	
	"stock_broker_application/src/app/watchlist/models"
	"errors"
	"fmt"
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

		defer func() {
			duration := time.Since(start)
			log.Printf("Completed in %v", duration)
		}()

		redisClient,_, err := utils.GetRedisClient(c,true)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.RedisConnectionError,
				},
				Error: constants.OperationFailedError,
			})
			return
		}

		authHeader := c.GetHeader(genericConstants.Authorization)

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Header,
					ErrorMessage: constants.ErrHeaderMissing,
				},
				Error: constants.OperationFailedError,
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, constants.Bearer)

		cacheKey := fmt.Sprintf(constants.BlacklistedCacheKey, tokenString)

		exists, err := redisClient.Exists(c, cacheKey).Result()

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.RedisConnectionError,
				},
				Error: constants.OperationFailedError,
			})
			return
		}

		if exists > 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Token,
					ErrorMessage: constants.InvalidTokenError,
				},
				Error: constants.OperationFailedError,
			})
			return
		}

		token, err := utils.ParseToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Token,
					ErrorMessage: err.Error(),
				},
				Error: constants.OperationFailedError,
			})
			return
		}

		claims, err := utils.VerifyToken(token)

		if err != nil {
			if errors.Is(err, genericConstants.ClaimMappingFailedError) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Token,
						ErrorMessage: genericConstants.ErrClaimMappingFailed,
					},
					Error: constants.OperationFailedError,
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Token,
						ErrorMessage: genericConstants.ErrTokenIsInvalid,
					},
					Error: constants.OperationFailedError,
				})
			return
		}

		username, ok := claims[constants.Subject].(string)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Token,
						ErrorMessage: constants.ErrUsernameNotFoundInJWT,
					},
					Error: constants.OperationFailedError,
				})
			return
		}


		c.Set(constants.User, username)

		c.Next()
	}
}
