package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"stock_broker_application/src/utils"
	"strings"
	"time"
	"watchlists/commons/constants"
	"watchlists/models"
	structModels "watchlists/models"
	"watchlists/repository"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

type WatchlistsService struct {
	watchlistsRepository repository.WatchlistsRepository // we have taken a field that is an interface in repository
}

func NewWatchlistsService(watchlistsRepository repository.WatchlistsRepository) *WatchlistsService {
	return &WatchlistsService{
		watchlistsRepository: watchlistsRepository, //constructor
	}
}

func (service *WatchlistsService) ADGtoWatchlist(ctx context.Context, spanCtx context.Context, bffWatchlistsRequest models.BFFAdgToWatchlistRequest, username string) ([]structModels.WatchlistWithId, []string, error) {
	postgresClient := utils.GetPostgresClient().GormDB
	redisClient, _ := utils.GetRedisClient()
	var warnings []string
	var cacheKey, cachedData string

	users, err := service.watchlistsRepository.GetUserId(spanCtx, postgresClient, username)
	if err != nil {
		if err.Error() == constants.UserNotFoundError {
			return nil, nil, errors.New(constants.UserNotFoundError)
		}
		return nil, nil, errors.New(constants.AuthenticationFailedError)
	}

	action := models.ActionType(strings.ToUpper(string(bffWatchlistsRequest.Action)))

	switch action {

	case models.ADD:

		exists, err := service.watchlistsRepository.CheckScripExists(spanCtx, postgresClient, bffWatchlistsRequest.ScripId)
		if err != nil {
			return nil, nil, fmt.Errorf(constants.QueryError, err)
		}
		if !exists {
			return nil, nil, errors.New(constants.ScripIdNotFoundError)
		}

		watchlistsDB, notValidIds, err := service.watchlistsRepository.GetValidWatchlists(
			spanCtx, postgresClient, users.ID, bffWatchlistsRequest.WatchlistIds,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(constants.QueryError, err)
		}

		if len(watchlistsDB) == 0 {
			return nil, nil, errors.New(constants.NoValidWatchlistIdsError)
		}

		if len(notValidIds) > 0 {
			warnings = append(warnings, fmt.Sprintf("Invalid watchlistIds skipped: %v", notValidIds))
		}

		var capNotFullIds []uint64
		for _, wl := range watchlistsDB {
			if wl.ScripCount >= 10 {
				warnings = append(warnings, fmt.Sprintf("Watchlist %d full (max 10)", wl.Id))
			}
			capNotFullIds = append(capNotFullIds, wl.Id)
		}

		if len(capNotFullIds) == 0 {
			return nil, warnings, nil
		}

		watchlistIdNames, err := service.watchlistsRepository.AddScripToWatchist(
			spanCtx,
			postgresClient,
			users.ID,
			capNotFullIds,
			bffWatchlistsRequest.ScripId,
		)
		if err != nil {
			if err.Error() == constants.UniqueConstraintViolationError {
				var duplicateIds []uint64

				addedMap := make(map[uint64]struct{})
				for _, wl := range watchlistIdNames {
					addedMap[uint64(wl.Watchlist_ID)] = struct{}{}
				}

				for _, id := range capNotFullIds {
					if _, ok := addedMap[id]; !ok {
						duplicateIds = append(duplicateIds, id)
					}
				}

				if len(duplicateIds) > 0 {
					warnings = append(warnings,
						fmt.Sprintf("Scrip %s already exists in watchlist: %v",
							bffWatchlistsRequest.ScripId,
							duplicateIds,
						),
					)
				}
				return watchlistIdNames, warnings, nil
			}
			return nil, nil, err
		}

		return watchlistIdNames, warnings, nil

	case models.DEL:

		watchlists, err := service.watchlistsRepository.DeleteScripFromWatchlist(
			spanCtx,
			postgresClient,
			users.ID,
			bffWatchlistsRequest.WatchlistIds,
			bffWatchlistsRequest.ScripId,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(constants.QueryError, err)
		}

		if len(watchlists) == 0 {
			return nil, warnings, nil
		}

		return watchlists, warnings, nil

	case models.GET:
		if strings.TrimSpace(bffWatchlistsRequest.ScripId) == "" {
			return nil, nil, errors.New(constants.ScripIdNotFoundError)
		}

		if redisClient != nil {
			cacheKey = fmt.Sprintf("userId:%d:scripId:%s:watchlists", users.ID, bffWatchlistsRequest.ScripId)
			cachedData, _ = redisClient.Get(ctx, cacheKey).Result()

			if len(cachedData) > 0 {
				var chachedWatchlists []structModels.WatchlistWithId
				err = json.Unmarshal([]byte(cachedData), &chachedWatchlists)
				if err != nil {
					log.Error("Error while Unmarshalling", zap.Error(err))
				}
				log.Info("Returning Data from Redis Cache")
				return chachedWatchlists, nil, nil
			}
		}
		log.Info("Failed to get data from redis, getting data redy from database...")

		watchlists, err := service.watchlistsRepository.GetUserWatchlists(
			spanCtx, postgresClient, users.ID, bffWatchlistsRequest.ScripId,
		)
		if err != nil {
			return nil, nil, fmt.Errorf(constants.QueryError, err)
		}

		if len(watchlists) == 0 {
			return nil, nil, errors.New(constants.WatchlistNotFoundError)
		}

		cacheValue, err := json.Marshal(watchlists)
		if err != nil {
			log.Error("Error While Marshalling", zap.Error(err))
		}

		err = redisClient.Set(ctx, cacheKey, cacheValue, time.Minute*30).Err()
		if err != nil {
			log.Error("Error while setting data in cache")
		}
		return watchlists, nil, nil
	}
	return nil, nil, err
}
