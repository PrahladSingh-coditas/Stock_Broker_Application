package business

import (
	"context"
	"errors"
	"fmt"
	"time"

	genericConstants "stock_broker_application/src/constants"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type LogoutService struct {
	redisClient *redis.Client
}

func NewLogoutService(redisClient *redis.Client) *LogoutService {
	return &LogoutService{
		redisClient: redisClient,
	}
}

func (service *LogoutService) Logout(ctx context.Context, token string, ttl time.Duration) error {
	logrus.SetLevel(logrus.WarnLevel)
	if service.redisClient == nil {
		logrus.Error("redis client not initialized")
		return errors.New(genericConstants.RedisClientNotInitialized)
	}

	key := fmt.Sprintf("BLACKLISTED_TOKEN_%s", token)
	err := service.redisClient.Set(ctx, key, 1, ttl).Err()
	if err != nil {
		logger.Info("error saving key in redis", err)
		return err
	}
	logger.Infof("key: %s, set successfully in redis", key)
	return nil
}
