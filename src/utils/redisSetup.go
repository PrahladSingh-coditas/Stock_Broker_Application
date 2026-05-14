package utils

import (
	"context"
	"stock_broker_application/src/constants"
	"sync"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var once sync.Once
var redisErr error

var mockRedisClient *redis.Client
var mockRedisController redismock.ClientMock

func initRedisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     constants.Address,
		Password: constants.Password,
		DB:       constants.DB,
	})

	_, err := client.Ping(context.Background()).Result()

	RedisClient = client
	redisErr = err
}

func initMockRedisClient(t *testing.T) {
	redisClient, mockRedis := redismock.NewClientMock()
	t.Cleanup(func() { redisClient.Close() })
	mockRedisClient = redisClient
	mockRedisController = mockRedis
}

func GetRedisClient(isMock bool) (*redis.Client, redismock.ClientMock, error) {
	if !isMock {
		if RedisClient != nil {
			return RedisClient, nil, nil
		}

		once.Do(func() {
			initRedisClient()
		})

		return RedisClient, nil, redisErr
	} else {
		if mockRedisClient != nil {
			return mockRedisClient, mockRedisController, nil
		}

		once.Do(func() {
			initMockRedisClient(&testing.T{})
		})

		return mockRedisClient, mockRedisController, nil
	}
}
