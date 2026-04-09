package utils

import (
	"context"
	"stock_broker_application/src/constants"
	"sync"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisOnce sync.Once
var redisErr error

func initRedisConfig(ctx context.Context) {
	rc := redis.NewClient(&redis.Options{
		Addr:     constants.Localhost,
		Password: "",
		DB:       0,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return
	}
	redisClient = rc
	redisErr = err
}

func GetRedisClient(ctx context.Context) (*redis.Client, error) {
	if redisClient != nil {
		return redisClient, nil
	}
	redisOnce.Do(func() {
		initRedisConfig(ctx)
	})
	return redisClient, redisErr
}
