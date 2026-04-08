package utils

import (
	"context"
	"fmt"
	"stock_broker_application/src/constants"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedisConfig() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       constants.RedisDB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %v", err)
	}

	redisClient = rdb
	return nil
}

func GetRedisClient() *redis.Client {
	return redisClient
}
