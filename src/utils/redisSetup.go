package utils

import (
	"context"
	"fmt"
	"stock_broker_application/src/constants"
	"sync"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient *redis.Client

var once sync.Once

func InitRedisConfg() error {
	rc := redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddress,
		Password: constants.RedisPassword,
		DB:       constants.RedisDB,
		Protocol: constants.RedisProtocol,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf(constants.ErrRedisInitFailed)
	}
	fmt.Println(constants.RedisConnectionSuccess)

	RedisClient = rc
	return nil
}

func GetRedisClient() (*redis.Client, error) {
	err := InitRedisConfg()
	if err != nil {
		return nil, err
	}
	return RedisClient, nil
}
