package utils

import (
	"context"
	"fmt"
	"log"
	"stock_broker_application/src/constants"
	"sync"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient *redis.Client

var once sync.Once

func InitRedisConfg() error {
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Error connecting to Redis")
	}
	fmt.Println("Connected to Redis!")

	RedisClient = rc
	return nil
}

func GetRedisClient() *redis.Client {
	err := InitRedisConfg()
	if err != nil {
		log.Fatalf(constants.ErrRedisInitFailed)
	}
	return RedisClient
}

// func SetRedisData(ctx context.Context, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, error) {
// 	err := client.Set(ctx, "greeting", "Hello, Redis!", 0).Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func GetRedisData(ctx context.Context, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, error){
// 	value, err := client.Get(ctx, "greeting").Result()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
