package middleware

import (
	"authentication/commons/constants"
	"authentication/models"
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		defer func() {
			duration := time.Since(start)
			log.Printf("Completed in %v", duration)
		}()

		redisClient, err := utils.GetRedisClient(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.RedisConnectionError,
				},
				Error: constants.OperationFailed,
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
				Error: constants.OperationFailed,
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
				Error: constants.OperationFailed,
			})
			return
		}

		if exists > 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Token,
					ErrorMessage: constants.InvalidTokenError,
				},
				Error: constants.OperationFailed,
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
				Error: constants.OperationFailed,
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
					Error: constants.OperationFailed,
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Token,
						ErrorMessage: genericConstants.ErrTokenIsInvalid,
					},
					Error: constants.OperationFailed,
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
					Error: constants.OperationFailed,
				})
			return
		}

		expiryTime, ok := claims[constants.Expiry].(float64)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Token,
						ErrorMessage: constants.ErrExpiryTimeNotFoundInJWT,
					},
					Error: constants.OperationFailed,
				})
			return
		}

		c.Set(constants.User, username)
		c.Set(constants.Token, tokenString)
		c.Set(constants.ExpiryTime, int64(expiryTime))

		c.Next()
	}
}
