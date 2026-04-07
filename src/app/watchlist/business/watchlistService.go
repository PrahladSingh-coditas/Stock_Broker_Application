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

	"stock_broker_application/src/utils"
)

type WatchlistService struct {
	watchlistRepository repository.WatchlistRepository
}

func NewWatchlistService(watchlistRepository repository.WatchlistRepository) *WatchlistService {
	return &WatchlistService{
		watchlistRepository: watchlistRepository,
	}
}

func (service *WatchlistService) Watchlist(ctx context.Context, spanCtx context.Context, username string, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, []string, error) {
	postgreClient := utils.GetPostgresClient()
	client := postgreClient.GormDB

	redisClient := utils.GetRedisClient()

	actionType := bffAdgToWatchlistRequest.Action

	//get user from db and to extract and use user.ID in ADG operations
	user, err := service.watchlistRepository.GetUserFromDb(spanCtx, client, username)
	if err != nil {
		return nil, nil, errors.New(constants.ErrUserNotFound)
	}

	switch actionType {
	case models.ADD:

		addedWatchlists, warningsResult, err := service.watchlistRepository.WatchlistAddOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
		if err != nil {
			return nil, nil, err
		}

		var warningsAll []string
		for _, w := range warningsResult {
			msg := "Watchlist " + strconv.FormatUint(w.WatchlistId, 10) + ": " + w.Reason
			warningsAll = append(warningsAll, msg)
		}

		var warningsNotOfUser []string //optimisation as recommended
		for _, val := range warningsAll {
			if strings.Contains(val, constants.ErrWatchlistsNotOfUser) {
				warningsNotOfUser = append(warningsNotOfUser, val)
			} else {
				continue
			}
		}

		if len(warningsNotOfUser) == len(bffAdgToWatchlistRequest.WatchlistIds) {
			return nil, nil, errors.New(constants.ErrWatchlistsNotOfUser)
		} else {
			return addedWatchlists, warningsAll, nil
		}

	case models.DEL:

		deletedWatchlists, warningsResult, err := service.watchlistRepository.WatchlistDeleteOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
		if err != nil {
			return nil, nil, err
		}

		var warningsAll []string
		for _, w := range warningsResult {
			msg := "Watchlist " + strconv.FormatUint(w.WatchlistId, 10) + ": " + w.Reason
			warningsAll = append(warningsAll, msg)
		}

		fmt.Println("warnings:", warningsAll)
		fmt.Println("deleted:", deletedWatchlists)

		var warningsNotOfUser []string
		for _, val := range warningsAll {
			if strings.Contains(val, constants.ErrWatchlistsNotOfUser) {
				warningsNotOfUser = append(warningsNotOfUser, val)
			} else {
				continue
			}
		}

		if len(warningsNotOfUser) == len(bffAdgToWatchlistRequest.WatchlistIds) {
			return nil, nil, errors.New(constants.ErrWatchlistsNotOfUser)
		} else {
			return deletedWatchlists, warningsAll, nil
		}

	case models.GET:
		key := fmt.Sprintf("Watchlist of:%d ScripId:%s", user.ID, bffAdgToWatchlistRequest.ScripId)

		if redisClient == nil {
			fmt.Println("Redisclient is nil")
		}

		data, err := redisClient.Get(ctx, key).Result()
		if len(data) > 0 {
			var CachedWatchlistIds []models.WatchlistWithId

			err = json.Unmarshal([]byte(data), &CachedWatchlistIds)
			if err != nil {
				return nil, nil, err
			}

			return CachedWatchlistIds, nil, nil
		} else {
			WatchlistNamewithId, err := service.watchlistRepository.WatchlistGetOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
			if err != nil {
				return nil, nil, errors.New(constants.ErrNoRowsAffected)
			}

			if len(WatchlistNamewithId) == 0 {
				return nil, nil, errors.New(constants.ErrWatchlistNotFound)
			}

			jsonData, err := json.Marshal(WatchlistNamewithId)
			if err != nil {
				return nil, nil, err
			}

			err = redisClient.Set(ctx, key, jsonData, 30*time.Minute).Err()
			if err != nil {
				return nil, nil, err
			}

			var warnings []string
			warnings = append(warnings, constants.ErrWatchlistNotRequired)
			if len(bffAdgToWatchlistRequest.WatchlistIds) != 0 {
				return WatchlistNamewithId, warnings, nil
			}

			return WatchlistNamewithId, nil, nil

		}

	default:
		return nil, nil, errors.New(constants.ErrInvalidAction)

	}

}
