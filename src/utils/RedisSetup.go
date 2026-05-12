package utils

import (
	"context"
	"stock_broker_application/src/constants"
	"sync"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisOnce sync.Once
var redisErr error

var mockRedisClient *redis.Client
var mockRedisController redismock.ClientMock

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

func getMockRedisClient(t *testing.T) {
	client, mock := redismock.NewClientMock()
	t.Cleanup(func() { client.Close() })
	mockRedisClient = client
	mockRedisController = mock
}

func GetRedisClient(ctx context.Context, isReal bool) (*redis.Client, redismock.ClientMock, error) {
	if isReal {
		if redisClient != nil {
			return redisClient, nil, nil
		}
		redisOnce.Do(func() {
			initRedisConfig(ctx)
		})
		return redisClient, nil, redisErr
	} else {
		if mockRedisClient != nil {
			return mockRedisClient, mockRedisController, nil
		}
		redisOnce.Do(func() {
			getMockRedisClient(&testing.T{})
		})
		return mockRedisClient, mockRedisController, nil
	}
}
