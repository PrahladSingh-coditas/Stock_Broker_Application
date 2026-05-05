package business

import (
	"authentication/commons/constants"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type LogoutUserService struct {
	redisClient *redis.Client
	redisError  error
}

func NewLogoutUser(redisClient *redis.Client, redisError error) *LogoutUserService {
	return &LogoutUserService{
		redisClient: redisClient,
		redisError:  redisError,
	}
}
func (service *LogoutUserService) LogoutUser(ctx context.Context, balckListedToken string, expiry int64) error {
	cacheKey := fmt.Sprintf("BLACKLISTE_TOKEN:%s", balckListedToken)
	cacheValue := 1

	if service.redisError != nil {
		fmt.Println("RedisConnectionFailedError")
		return errors.New(constants.RedisConnectionFailedError)
	}

	duration := expiry - time.Now().Unix()
	timeToLive := time.Duration(duration) * time.Second
	if service.redisClient != nil {
		err := service.redisClient.Set(ctx, cacheKey, cacheValue, timeToLive).Err()
		if err != nil {
			return errors.New(constants.RedisOperationFailedError)
		} else {
			return nil
		}
	}
	return nil
}
