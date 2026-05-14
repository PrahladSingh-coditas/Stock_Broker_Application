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
}

func NewLogoutUser(redisClient *redis.Client) *LogoutUserService {
	return &LogoutUserService{
		redisClient: redisClient,
	}
}
func (service *LogoutUserService) LogoutUser(ctx context.Context, balckListedToken string, expiry int64) error {
	cacheKey := fmt.Sprintf("BLACKLISTE_TOKEN:%s", balckListedToken)
	cacheValue := 1

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
