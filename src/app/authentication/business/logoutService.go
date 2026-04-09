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
func (service *LogoutUserService) LogoutUser(ctx context.Context, spanCtx context.Context, balckListedToken string, timeToLeave time.Duration) error {
	cacheKey := fmt.Sprintf("BLACKLISTE_TOKEN:%s", balckListedToken)
	cacheValue := 1

	redisClient, redisError := utils.GetRedisClient()
	if redisError != nil {
		fmt.Println("RedisConnectionFailedError")
		return errors.New(constants.RedisConnectionFailedError)
	}

	if redisClient != nil {
		err := redisClient.Set(ctx, cacheKey, cacheValue, timeToLeave).Err()
		if err != nil {
			return errors.New(constants.RedisSetOperationError)
		} else {
			return nil
		}
	}
	return nil
}
