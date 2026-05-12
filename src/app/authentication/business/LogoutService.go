package business

import (
	"authentication/commons/constants"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type LogoutUserService struct {
	redisClient *redis.Client
}

func NewLogoutUserService(redisClient *redis.Client) *LogoutUserService {
	return &LogoutUserService{redisClient: redisClient}
}

func (service *LogoutUserService) LogoutUser(ctx context.Context, logger *logrus.Logger, tokenString string, expiryTime int64) error {

	duration := expiryTime - time.Now().Unix()
	ttl := time.Duration(duration) * time.Second

	cacheKey := fmt.Sprintf(constants.BlacklistedCacheKey, tokenString)

	err := service.redisClient.Set(ctx, cacheKey, 1, ttl).Err()
	
	if err != nil {
		logger.Error(constants.RedisSetOperationError)
		return errors.New(constants.RedisSetOperationError)
	}

	return nil
}
