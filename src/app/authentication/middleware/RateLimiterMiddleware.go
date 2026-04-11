package middleware

import (
	"authentication/commons/constants"
	"authentication/models"
	"encoding/json"
	"fmt"
	"net/http"
	"stock_broker_application/src/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var bucket *models.TokenBucket
		username := ctx.GetString(constants.User)
		redisClient, err := utils.GetRedisClient(ctx)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.RedisConnectionError,
				},
				Error: constants.OperationFailed,
			})
			return
		}

		cacheKey := fmt.Sprintf("RATE_LIMIT_USERNAME_%s", username)

		data, err := redisClient.Get(ctx, cacheKey).Result()

		if len(data) == 0 {
			bucket = models.GetNewTokenBucket(3, 3)

			result := bucket.Allow()

			marshalledData, err := json.Marshal(&bucket)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Redis,
						ErrorMessage: "error in marshalling",
					},
					Error: constants.OperationFailed,
				})
				return
			}

			err = redisClient.Set(ctx, cacheKey, marshalledData, 1*time.Hour).Err()

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Redis,
						ErrorMessage: "error in adding data in redis",
					},
					Error: constants.OperationFailed,
				})
				return
			}

			if !result {
				ctx.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorAPIResponse{
					Error: "Too many request",
				})
				return
			}
		} else {
			err = json.Unmarshal([]byte(data), &bucket)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Redis,
						ErrorMessage: "error in unmarshalling",
					},
					Error: constants.OperationFailed,
				})
				return
			}

			result := bucket.Allow()

			marshalledData, err := json.Marshal(&bucket)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Redis,
						ErrorMessage: "error in marshalling",
					},
					Error: constants.OperationFailed,
				})
				return
			}

			err = redisClient.Set(ctx, cacheKey, marshalledData, 1*time.Hour).Err()

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
					Message: models.ErrorMessage{
						Key:          constants.Redis,
						ErrorMessage: "error in adding data in redis",
					},
					Error: constants.OperationFailed,
				})
				return
			}

			if !result {
				ctx.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorAPIResponse{
					Error: "Too many request",
				})
				return
			}
		}

		ctx.Next()
	}
}
