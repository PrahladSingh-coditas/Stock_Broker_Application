package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"watchlists/commons/constants"
	"watchlists/models"
	"watchlists/repository"

	"github.com/pingcap/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdgStoWatchlistService struct {
	adgStoWatchlistRepository repository.AdgStoWatchlistRepository
	rdb                       *redis.Client
}

func NewadgStoWatchlistService(repo repository.AdgStoWatchlistRepository, rdb *redis.Client) *AdgStoWatchlistService {
	return &AdgStoWatchlistService{
		adgStoWatchlistRepository: repo,
		rdb:                       rdb,
	}
}

func (service *AdgStoWatchlistService) AdgStoWatchlist(ctx context.Context, username string, request models.BffAdgStoWatchlistRequest) ([]string, []models.WatchlistWithId, error) {
	// rdb := utils.GetRedisClient()
	ReqAction := models.Actiontype(strings.ToUpper(string(request.Action)))

	user, err := service.adgStoWatchlistRepository.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, constants.UserNotFoundError
		}
		return nil, nil, err
	}

	switch ReqAction {
	case models.AddAction:
		warnings := []string{}
		respWatchlistWithIds := []models.WatchlistWithId{}

		validIds, insertedIds, err := service.adgStoWatchlistRepository.AddScripToWatchlists(ctx, user.ID, request.ScripId, request.WatchlistIds)
		if err != nil {
			return warnings, respWatchlistWithIds, err
		}

		//no valid watchlists in the request - error
		if len(validIds) == 0 {
			return warnings, respWatchlistWithIds, constants.InvalidWatchlistsError
		}

		if len(insertedIds) == 0 {
			return warnings, respWatchlistWithIds, constants.ScripNotAddedToAnyWatchlistError
		}

		//list of invalid watchlistIds
		invalidIds := []uint64{}
		validMap := make(map[uint64]bool)
		for _, id := range validIds {
			validMap[id] = true
		}
		for _, id := range request.WatchlistIds {
			if !validMap[id] {
				invalidIds = append(invalidIds, id)
			}
		}
		if len(invalidIds) > 0 {
			warnings = append(warnings, fmt.Sprintf("watchlistIds %v does not belong to the user", invalidIds))
		}

		//watchlistIds taht already had the scrip or were full
		skippedIds := []uint64{}
		skippedMap := make(map[uint64]bool)
		for _, id := range insertedIds {
			skippedMap[id] = true
		}
		for _, id := range validIds {
			if !skippedMap[id] {
				skippedIds = append(skippedIds, id)
			}
		}
		if len(skippedIds) > 0 {
			warnings = append(warnings, fmt.Sprintf("scrip not added to watchlistIds %v, either full or duplicate", skippedIds))
		}

		resp, err := service.adgStoWatchlistRepository.GetWatchlistDetails(ctx, insertedIds)
		if err != nil {
			return warnings, respWatchlistWithIds, constants.DatabaseQueryError
		}
		for _, r := range resp {
			respWatchlistWithIds = append(respWatchlistWithIds, models.WatchlistWithId{
				Id:   r.Id,
				Name: r.WatchlistName,
			})
		}

		key := fmt.Sprintf("watchlists:%d:%s", user.ID, request.ScripId)
		err = service.rdb.Del(ctx, key).Err()
		if err != nil {
			log.Error("error deleting cached data: ", zap.Error(err))
		}
		return warnings, respWatchlistWithIds, nil

	case models.DelAction:
		warnings := []string{}
		respWatchlistWithIds := []models.WatchlistWithId{}

		validWatchlists, err := service.adgStoWatchlistRepository.DelScripFromWatchlists(ctx, user.ID, request.ScripId, request.WatchlistIds)
		if err != nil {
			return warnings, respWatchlistWithIds, err
		}
		if len(validWatchlists) == 0 {
			return warnings, respWatchlistWithIds, constants.InvalidWatchlistsError
		}

		//check for invalid ids for warnings
		usersIdsMap := make(map[uint64]bool)
		inValidIds := []uint64{}

		for _, id := range validWatchlists {
			usersIdsMap[id] = true
		}

		for _, id := range request.WatchlistIds {
			if !usersIdsMap[id] {
				inValidIds = append(inValidIds, id)
			}
		}

		if len(inValidIds) > 0 {
			warnings = append(warnings, fmt.Sprintf("watchlistIds %v do not belong to the user", inValidIds))
		}

		response, err := service.adgStoWatchlistRepository.GetWatchlistDetails(ctx, validWatchlists)
		if err != nil {
			return warnings, respWatchlistWithIds, err
		}
		for _, w := range response {
			respWatchlistWithIds = append(respWatchlistWithIds, models.WatchlistWithId{
				Id:   w.Id,
				Name: w.WatchlistName,
			})
		}

		key := fmt.Sprintf("watchlists:%d:%s", user.ID, request.ScripId)
		err = service.rdb.Del(ctx, key).Err()
		if err != nil {
			log.Error(constants.ErrDeletingFromCache, zap.Error(err))
		}

		return warnings, respWatchlistWithIds, nil

	case models.GetAction:
		warnings := []string{}
		respWatchlistWithId := []models.WatchlistWithId{}

		key := fmt.Sprintf("watchlists:%d:%s", user.ID, request.ScripId)
		val, err := service.rdb.Get(ctx, key).Result()
		if err == nil {
			var cachedIds []models.WatchlistWithId
			if err := json.Unmarshal([]byte(val), &cachedIds); err == nil {
				if len(cachedIds) == 0 {
					return warnings, cachedIds, constants.ScripNotInWatchlistsError
				}
				return warnings, cachedIds, nil
			}
			log.Error(constants.ErrUnmarshallingCache, zap.Error(err))
		}

		if err != nil && err != redis.Nil {
			log.Error("redis error: ", zap.Error(err))
		}
		//not in redis, going to db
		watchlists, err := service.adgStoWatchlistRepository.GetScripFromWatchlists(ctx, user.ID, request.ScripId)
		if err != nil {
			return warnings, respWatchlistWithId, err
		}
		if len(watchlists) == 0 {
			return warnings, respWatchlistWithId, constants.ScripNotInWatchlistsError
		}

		//if fetched, set key-value to redis
		data, err := json.Marshal(watchlists)
		if err == nil {
			if err = service.rdb.Set(ctx, key, data, constants.RedisKeyTTL).Err(); err != nil {
				log.Error(constants.ErrSavingToCache, zap.Error(err))
			}
		} else {
			log.Error(constants.ErrMarshallingCache, zap.Error(err))
		}

		return warnings, watchlists, nil

	default:
		return nil, nil, constants.InvalidActionTypeError
	}
}
