package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"watchlist/commons/constants"
	"watchlist/models"
	"watchlist/repository"

	"github.com/pingcap/log"
	"go.uber.org/zap"

	"stock_broker_application/src/utils"
)

type WatchlistService struct {
	watchlistRepository repository.WatchlistRepositoryInterface
}

func NewWatchlistService(watchlistRepository repository.WatchlistRepositoryInterface) *WatchlistService {
	return &WatchlistService{
		watchlistRepository: watchlistRepository,
	}
}

func (service *WatchlistService) Watchlist(ctx context.Context, spanCtx context.Context, username string, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, []string, error) {
	postgreClient := utils.GetPostgresClient()
	client := postgreClient.GormDB

	//get user from db and to extract and use user.ID in ADG operations
	user, err := service.watchlistRepository.GetUserFromDb(spanCtx, client, username)
	if err != nil {
		return nil, nil, errors.New(constants.ErrUserNotFound)
	}

	//redis connection and define key for redis
	redisClient := utils.GetRedisClient()
	redisKey := fmt.Sprintf("Watchlist of:%d ScripId:%s", user.ID, bffAdgToWatchlistRequest.ScripId)

	//extract action from request and switch case on it
	actionType := bffAdgToWatchlistRequest.Action

	switch actionType {
	case models.ADD:

		addedWatchlists, warningsResult, err := service.watchlistRepository.WatchlistAddOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
		if err != nil {
			return nil, nil, err
		} else if redisClient != nil {
			err := redisClient.Del(ctx, redisKey).Err()
			if err != nil {
				log.Error("Failed to invalidate cache after ADD", zap.Error(err))
			} else {
				log.Info("Cache invalidated after ADD")
			}
		}

		var warningsAll []string
		for _, w := range warningsResult {
			msg := "Watchlist " + strconv.FormatUint(w.WatchlistId, 10) + ": " + w.Reason
			warningsAll = append(warningsAll, msg)
		}

		var countWatchlistNotOfUser int 
		for _, val := range warningsAll {
			if strings.Contains(val, constants.ErrWatchlistsNotOfUser) {
				countWatchlistNotOfUser++
			} else {
				continue
			}
		}

		if countWatchlistNotOfUser == len(bffAdgToWatchlistRequest.WatchlistIds) {
			return nil, nil, errors.New(constants.ErrWatchlistsNotOfUser)
		} else {
			return addedWatchlists, warningsAll, nil
		}

	case models.DEL:

		deletedWatchlists, warningsResult, err := service.watchlistRepository.WatchlistDeleteOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
		if err != nil {
			return nil, nil, err
		} else if redisClient != nil {
			err := redisClient.Del(ctx, redisKey).Err()
			if err != nil {
				log.Error("Failed to invalidate cache after DEL", zap.Error(err))
			} else {
				log.Info("Cache invalidated after DEL")
			}
		}

		var warningsAll []string
		for _, w := range warningsResult {
			msg := "Watchlist " + strconv.FormatUint(w.WatchlistId, 10) + ": " + w.Reason
			warningsAll = append(warningsAll, msg)
		}

		fmt.Println("warnings:", warningsAll)
		fmt.Println("deleted:", deletedWatchlists)

		var countWatchlistNotOfUser int
		for _, val := range warningsAll {
			if strings.Contains(val, constants.ErrWatchlistsNotOfUser) {
				countWatchlistNotOfUser++
			} else {
				continue
			}
		}

		if countWatchlistNotOfUser == len(bffAdgToWatchlistRequest.WatchlistIds) {
			return nil, nil, errors.New(constants.ErrWatchlistsNotOfUser)
		} else {
			return deletedWatchlists, warningsAll, nil
		}

	case models.GET:
		var CachedWatchlistIds []models.WatchlistWithId
		//only warning: watchlists not required for get
		var warnings []string

		if redisClient != nil {
			data, err := redisClient.Get(ctx, redisKey).Result()
			if err == nil && len(data) > 0 {
				err := json.Unmarshal([]byte(data), &CachedWatchlistIds)
				if err == nil {
					fmt.Println("\nData fetched from Redis!!!!!!!!!!!!!!!!!!!!!!")

					warnings = append(warnings, constants.ErrWatchlistNotRequired)
					if len(bffAdgToWatchlistRequest.WatchlistIds) != 0 {
						return CachedWatchlistIds, warnings, nil
					}

					return CachedWatchlistIds, nil, nil
				} else {
					return nil, nil, errors.New(constants.ErrJsonUnmarshalFailed)
				}
			} 
		}

		WatchlistNamewithId, err := service.watchlistRepository.WatchlistGetOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
		if err != nil {
			return nil, nil, errors.New(constants.ErrNoRowsAffected)
		}

		if len(WatchlistNamewithId) == 0 {
			return nil, nil, errors.New(constants.ErrWatchlistNotFound)
		}

		if redisClient != nil {
			jsonData, err := json.Marshal(WatchlistNamewithId)
			if err == nil {
				_ = redisClient.Set(ctx, redisKey, jsonData, 30*time.Minute).Err()	
			} else {
				return nil,nil,errors.New(constants.ErrJsonMarshalFailed)
			}
		}

		fmt.Println("\nData fetched from DB!!!!!!!!!!!!!!!!!!!!!!")
		warnings = append(warnings, constants.ErrWatchlistNotRequired)
		if len(bffAdgToWatchlistRequest.WatchlistIds) != 0 {
			return WatchlistNamewithId, warnings, nil
		}

		return WatchlistNamewithId, nil, nil

	default:
		return nil, nil, errors.New(constants.ErrInvalidAction)

	}

}
