package business

import (
	"context"
	"errors"
	"fmt"
	"stock_broker_application/src/utils"
	"strings"
	"watchlists/commons/constants"
	"watchlists/models"
	structModels "watchlists/models"
	"watchlists/repository"
)

type WatchlistsService struct {
	watchlistsRepository repository.WatchlistsRepository // we have taken a field that is an interface in repository
}

func NewWatchlistsService(watchlistsRepository repository.WatchlistsRepository) *WatchlistsService {
	return &WatchlistsService{
		watchlistsRepository: watchlistsRepository, //constructor
	}
}

func (service *WatchlistsService) Watchlists(ctx context.Context, spanCtx context.Context, bffWatchlistsRequest models.BFFAdgToWatchlistRequest, username string) ([]structModels.WatchlistWithId, []string, error) {
	postgresClinet := utils.GetPostgresClient()
	client := postgresClinet.GormDB

	// var watchlistIdNames []structModels.WatchlistWithId
	var warnings []string
	users, err := service.watchlistsRepository.GetUserId(spanCtx, client, username)
	if err != nil {
		if err.Error() == constants.UserNotFoundError {
			return nil, nil, errors.New(constants.UserNotFoundError)
		}
		return nil, nil, errors.New(constants.AuthenticationFailedError)
	}

	action := models.ActionType(strings.ToUpper(string(bffWatchlistsRequest.Action)))

	switch action {

	case models.ADD:

		exists, err := service.watchlistsRepository.CheckScripExists(spanCtx, client, bffWatchlistsRequest.ScripId)
		if err != nil {
			return nil, nil, errors.New(constants.QueryError)
		}
		if !exists {
			return nil, nil, errors.New(constants.ScripIdNotFoundError)
		}

		watchlistsDB, notValidIds, err := service.watchlistsRepository.GetValidWatchlists(
			spanCtx, client, users.ID, bffWatchlistsRequest.WatchlistIds,
		)
		if err != nil {
			return nil, nil, errors.New(constants.QueryError)
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

		watchlistIdNames, err := service.watchlistsRepository.AddScripWithCTE(
			spanCtx,
			client,
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

		watchlists, err := service.watchlistsRepository.DeleteScripWithCTE(
			spanCtx,
			client,
			users.ID,
			bffWatchlistsRequest.WatchlistIds,
			bffWatchlistsRequest.ScripId,
		)
		if err != nil {
			return nil, nil, errors.New(constants.QueryError)
		}

		if len(watchlists) == 0 {
			return nil, warnings, nil
		}

		return watchlists, warnings, nil

	case models.GET:
		if strings.TrimSpace(bffWatchlistsRequest.ScripId) == "" {
			return nil, nil, errors.New(constants.ScripIdNotFoundError)
		}

		watchlists, err := service.watchlistsRepository.GetUserWatchlists(
			spanCtx, client, users.ID, bffWatchlistsRequest.ScripId,
		)
		if err != nil {
			return nil, nil, errors.New(constants.QueryError)
		}

		if len(watchlists) == 0 {
			return nil, nil, errors.New(constants.WatchlistNotFoundError)
		}

		return watchlists, nil, nil
	}
	return nil, nil, err
}
