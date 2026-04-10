package business

import (
	"authentication/commons/constants"
	"context"
	"errors"
	"fmt"
	"stock_broker_application/src/utils"
	"time"
)

type LogoutUserService struct{}

func NewLogoutUser() *LogoutUserService {
	return &LogoutUserService{}
}
func (service *LogoutUserService) LogoutUser(ctx context.Context, balckListedToken string, expiry int64) error {
	cacheKey := fmt.Sprintf("BLACKLISTE_TOKEN:%s", balckListedToken)
	cacheValue := 1

	redisClient, redisError := utils.GetRedisClient()
	if redisError != nil {
		fmt.Println("RedisConnectionFailedError")
		return errors.New(constants.RedisConnectionFailedError)
	}

	duration := expiry - time.Now().Unix()
	timeToLive := time.Duration(duration) * time.Second
	if redisClient != nil {
		err := redisClient.Set(ctx, cacheKey, cacheValue, timeToLive).Err()
		if err != nil {
			return errors.New(constants.RedisOperationFailedError)
		} else {
			return nil
		}
	}
	return nil
}
