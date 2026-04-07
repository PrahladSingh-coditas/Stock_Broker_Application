package utils

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var once sync.Once
var redisErr error

func initRedisClient() {
	rc := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err := rc.Ping(context.Background()).Result()

	RedisClient = rc
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
