package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"stock_broker_application/src/app/watchlist/commons/constants"
	"stock_broker_application/src/app/watchlist/models"
	"stock_broker_application/src/app/watchlist/repository"
	"stock_broker_application/src/utils"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// struct declaration
type WatchlistService struct {
	watchlistRepository repository.WatchlistRepository
}

// struct initialisation
func NewWatchlistService(watchlistRepository repository.WatchlistRepository) *WatchlistService {
	return &WatchlistService{
		watchlistRepository: watchlistRepository, // service local var = router local var
	}
}

func (service *WatchlistService) ServiceWatchlist(ctx context.Context, spanCtx context.Context, logger *logrus.Logger, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest, username string) ([]models.WatchlistWithID, []string, error) {
	redisClient,_, err := utils.GetRedisClient(ctx,true)

	if err != nil || redisClient == nil {
		logger.Error(constants.RedisConnectionError)
	}

	postgresClient := utils.GetPostgresClient()
	tx := postgresClient.GormDB.Begin()
	var watchlistsWithId []models.WatchlistWithID
	var warnings []string
	var response string

	userId, err := service.watchlistRepository.GetUserIdByUsername(ctx, tx, logger, username)

	if err != nil {
		return nil, nil, err
	}

	cacheKey := fmt.Sprintf(constants.UserIdScripIdKey, strconv.Itoa(int(*userId)), bffAdgToWatchlistRequest.ScripId)

	switch bffAdgToWatchlistRequest.Action {
	case models.GET:

		if redisClient != nil {
			response, err = redisClient.Get(ctx, cacheKey).Result()
		}

		if err == nil && len(response) != 0 {
			err := json.Unmarshal([]byte(response), &watchlistsWithId)
			if err != nil {
				logger.Error(constants.RedisUnmarshallingError)
			} else {
				logger.Info(constants.DataFetchFromRedis)
				break
			}

		}

		resultListsForGET, err := service.watchlistRepository.GetWatchlistsWithId(ctx, tx, logger, *userId, bffAdgToWatchlistRequest.ScripId)
		if err != nil {
			return nil, nil, errors.New(constants.ErrNoWatchlistForScripMsg)
		}

		if resultListsForGET.ScripCheckCount == 0 {
			return nil, nil, errors.New(constants.ErrScripNotFoundMsg)
		}

		if len(resultListsForGET.WatchlistId) == 0 {
			return nil, nil, errors.New(constants.ErrNoWatchlistForScripMsg)
		}

		for i := 0; i < len(resultListsForGET.WatchlistId); i++ {
			watchlistsWithId = append(watchlistsWithId, models.WatchlistWithID{
				WatchlistId:   uint64(resultListsForGET.WatchlistId[i]),
				WatchlistName: resultListsForGET.WatchlistName[i],
			})
		}

		marshalledData, err := json.Marshal(watchlistsWithId)
		if err != nil {
			logger.Error(constants.RedismarshallingError)
		}

		err = redisClient.Set(ctx, cacheKey, marshalledData, time.Minute*60).Err()
		if err != nil {
			logger.Error(constants.RedisDataAdditionError)
		}

	case models.DEL:

		if len(bffAdgToWatchlistRequest.WatchlistIds) == 0 {
			return nil, nil, errors.New(constants.EmptyWatchlistIdsError)
		}

		resultListsForDEL, err := service.watchlistRepository.DeleteScripsFromWatchlists(ctx, tx, logger, *userId, bffAdgToWatchlistRequest.WatchlistIds, bffAdgToWatchlistRequest.ScripId)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}

		for _, id := range resultListsForDEL.DeletedWatchlistIds {
			watchlistsWithId = append(watchlistsWithId, models.WatchlistWithID{
				WatchlistId: uint64(id),
			})
		}

		if len(resultListsForDEL.ValidWatchlistIds) == 0 {
			return nil, nil, errors.New(constants.ErrNoWatchlistForScripMsg)
		} else if len(resultListsForDEL.ValidWatchlistIds) != len(bffAdgToWatchlistRequest.WatchlistIds) {
			warnings = append(warnings, constants.InvalidWatchlistIdsError)
		}

		err = redisClient.Del(ctx, cacheKey).Err()

		if err != nil {
			logger.Error(constants.RedisDataDeletionError)
		}

	case models.ADD:

		if len(bffAdgToWatchlistRequest.WatchlistIds) == 0 {
			return nil, nil, errors.New(constants.EmptyWatchlistIdsError)
		}

		result, err := service.watchlistRepository.AddScripsToWatchlists(ctx, tx, logger, *userId, bffAdgToWatchlistRequest.WatchlistIds, bffAdgToWatchlistRequest.ScripId)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}

		if result.ScripCount == 0 {
			return nil, nil, errors.New(constants.ErrScripNotFoundMsg)
		}

		if len(result.AddedWatchlistIds) == 0 {
			warnings = append(warnings, constants.NoScripAddedToWatchlistMsg)
		}

		if len(result.SkippedWatchlistIds) > 0 {
			warnings = append(warnings, constants.AlreadyExistsInSomeWatchlistMsg)
		}

		if len(result.LimitExceededWatchlistIds) > 0 {
			warnings = append(warnings, constants.SomeWatchlistsReachedMaxLimitMsg)
		}

		if len(result.AddedWatchlistIds) != 0 {
			for _, id := range result.AddedWatchlistIds {
				watchlistsWithId = append(watchlistsWithId, models.WatchlistWithID{
					WatchlistId: uint64(id),
				})
			}
		}

		err = redisClient.Del(ctx, cacheKey).Err()

		if err != nil {
			logger.Error(constants.RedisDataDeletionError)
		}

	}

	if len(warnings) == 0 {
		warnings = append(warnings, constants.NoWarningsMsg)
	}

	tx.Commit()

	return watchlistsWithId, warnings, nil
}
