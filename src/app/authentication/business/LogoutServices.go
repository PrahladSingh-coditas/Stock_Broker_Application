package business

import (
	authConstants "authentication/commons/constants"
	"context"
	"errors"

	"stock_broker_application/src/constants"
	"time"

	"github.com/redis/go-redis/v9"
)

type LogoutUserService struct {
	RedisClient *redis.Client
}

func NewLogoutUserService(redisClient *redis.Client) *LogoutUserService {
	return &LogoutUserService{
		RedisClient: redisClient,
	}
}

func (user *LogoutUserService) LogoutUser(ctx context.Context, tokenString string, tokenExpiry int64, redisKey string) error {
	//checking if redis connection has failed or not
	if user.RedisClient == nil {
		return errors.New(constants.ErrRedisInitFailed)
	}

	//rediskey will be the token itself and value is set to 1 as given
	// redisKey := fmt.Sprintf(authConstants.BlacklistedToken, tokenString)
	redisValue := 1

	currentTime := time.Now().Unix()
	remainingSeconds := tokenExpiry - currentTime

	timeToLive := time.Duration(remainingSeconds) * time.Second

	err := user.RedisClient.Set(ctx, redisKey, redisValue, timeToLive).Err()
	if err != nil {
		return errors.New(authConstants.ErrFailedToBlacklist)
	}

	return nil
}
