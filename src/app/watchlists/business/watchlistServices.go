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

	var watchlistIdNames []structModels.WatchlistWithId
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

		//unique check
		duplicates, notDuplicates, err := service.watchlistsRepository.CheckDuplicate(
			spanCtx, client, capNotFullIds, bffWatchlistsRequest.ScripId,
		)
		if err != nil {
			return nil, nil, errors.New(constants.QueryError)
		}
		if len(duplicates) > 0 {
			warnings = append(warnings, fmt.Sprintf("Already exists in watchlists: %v", duplicates))
		}

		if len(notDuplicates) == 0 {
			return nil, warnings, nil
		}

		addedIds, err := service.watchlistsRepository.InsertWatchlistScrip(
			spanCtx, client, notDuplicates, bffWatchlistsRequest.ScripId,
		)
		if err != nil {
			return nil, nil, errors.New(constants.QueryError)
		}

		watches, err := service.watchlistsRepository.GetWatchlistDetails(spanCtx, client, addedIds)
		if err != nil {
			return nil, nil, errors.New(constants.QueryError)
		}

		for _, w := range watches {
			watchlistIdNames = append(watchlistIdNames, structModels.WatchlistWithId{
				Watchlist_ID:   int64(w.Id),
				Watchlist_Name: w.WatchlistName,
			})
		}

		return watchlistIdNames, warnings, nil

	case models.DEL:
		watchlistsDB, _, err := service.watchlistsRepository.GetValidWatchlists(
			spanCtx, client, users.ID, bffWatchlistsRequest.WatchlistIds,
		)
		if err != nil {
			return nil, nil, errors.New(constants.QueryError)
		}
		// for _, id := range notValidIds {
		// 	warnings = append(warnings,
		// 		fmt.Sprintf("Watchlist with id %d does not belong to the user hence skipped", id))
		// }
		if len(watchlistsDB) == 0 {
			return nil, warnings, nil
		}

		var watchlists []structModels.WatchlistWithId
		for _, wl := range watchlistsDB {
			watchlists = append(watchlists, structModels.WatchlistWithId{
				Watchlist_ID:   int64(wl.Id),
				Watchlist_Name: wl.WatchlistName,
			})
			err := service.watchlistsRepository.DeleteWatchlistScrip(spanCtx, client, watchlists, bffWatchlistsRequest.ScripId)
			if err != nil {
				return nil, nil, errors.New(constants.QueryError)
			}
		}

		return watchlists, warnings, nil

	case models.GET:
		if strings.TrimSpace(bffWatchlistsRequest.ScripId) == "" {
			return nil, nil, errors.New("invalid scripId")
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
