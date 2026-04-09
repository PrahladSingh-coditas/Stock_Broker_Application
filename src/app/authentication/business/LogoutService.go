package business

import (
	"authentication/commons/constants"
	"context"
	"errors"
	"fmt"
	"stock_broker_application/src/utils"
	"time"

	"github.com/sirupsen/logrus"
)

type LogoutUserService struct{}

func NewLogoutUserService() *LogoutUserService {
	return &LogoutUserService{}
}

func (service *LogoutUserService) LogoutUser(ctx context.Context, logger *logrus.Logger, tokenString string, expiryTime int64) error {
	redisClient, err := utils.GetRedisClient(ctx)

	if err != nil {
		logger.Error(constants.RedisConnectionError)
		return errors.New(constants.RedisConnectionError)
	}

	duration := expiryTime - time.Now().Unix()
	ttl := time.Duration(duration) * time.Second

	cacheKey := fmt.Sprintf(constants.BlacklistedCacheKey, tokenString)

	err = redisClient.Set(ctx, cacheKey, 1, ttl).Err()

	if err != nil {
		logger.Error(constants.RedisConnectionError)
		return errors.New(constants.RedisConnectionError)
	}

	return nil
}
