package utils

import (
	"context"
	"stock_broker_application/src/constants"
	"sync"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var once sync.Once
var redisErr error

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

func GetRedisClient() (*redis.Client, error) {
	if RedisClient != nil {
		return RedisClient, nil
	}

	once.Do(func() {
		initRedisClient()
	})

	return RedisClient, redisErr
}
