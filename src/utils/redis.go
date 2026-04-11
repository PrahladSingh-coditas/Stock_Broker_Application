package utils

import (
	"context"
	"fmt"
	"stock_broker_application/src/constants"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func InitRedis() (*redis.Client, error) {
	var err error //because the Do func doesnt't return anything so we'll need to capture that error ouside that func to return in
	//once.Do creates a connectio pool only once
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     constants.RedisAddr,
			Password: constants.RedisPassword,
			DB:       constants.RedisDB,
		})

		if _, pingErr := client.Ping(context.Background()).Result(); err != nil {
			err = fmt.Errorf("failed to connect to redis: %v", pingErr)
			return
		}

		//success, assign the client to the global variable
		redisClient = client
	})

	//error
	return redisClient, err
}

func GetRedisClient() *redis.Client {
	return redisClient
}
