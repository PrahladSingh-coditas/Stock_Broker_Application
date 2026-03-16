package business

import (
	"context"
	"errors"
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

func (service *WatchlistsService) Watchlists(ctx context.Context, spanCtx context.Context, bffWatchlistsRequest models.BFFAdgToWatchlistRequest, username string) ([]structModels.WatchlistWithId, error) {
	postgresClinet := utils.GetPostgresClient()
	client := postgresClinet.GormDB
	users, err := service.watchlistsRepository.GetUserId(spanCtx, client, username)
	if err != nil {
		if err.Error() == constants.UserNotFoundError {
			return nil, errors.New(constants.UserNotFoundError)
		}
		return nil, errors.New(constants.AuthenticationFailedError)
	}

	// action := bffWatchlistsRequest.Action
	// if !models.ActionType(action).IsValid() {
	// 	fmt.Println(models.ActionType(action))
	// 	return nil, errors.New(constants.InvalidActionTypeError)
	// }

	action := models.ActionType(strings.ToUpper(string(bffWatchlistsRequest.Action)))

	switch action {

	case models.ADD:

		_, err := service.watchlistsRepository.CheckScripExists(spanCtx, client, bffWatchlistsRequest.ScripId)
		if err != nil {
			if err.Error() == constants.ScripIdNotFoundError {
				return nil, errors.New(constants.ScripIdNotFoundError)
			}
			return nil, errors.New(constants.AuthenticationFailedError)
		}

	// 	watchlistsDB, err := service.watchlistsRepository.GetValidWatchlists(spanCtx, client, users.ID, bffWatchlistsRequest.WatchlistIds)
	// 	if err != nil {
	// 		return nil, errors.New(constants.QueryError)
	// 	}

	// 	for _, wl := range watchlistsDB {

	// 		if wl.ScripCount >= 10 {
	// 			continue
	// 		}

	// 		duplicate, _ := service.watchlistsRepository.CheckDuplicate(spanCtx, client, wl.Id, bffWatchlistsRequest.ScripId)
	// 		if duplicate {
	// 			continue
	// 		}

	// 		err := service.watchlistsRepository.InsertWatchlistScrip(spanCtx, client, wl.Id, bffWatchlistsRequest.ScripId)
	// 		if err != nil {
	// 			return nil, errors.New(constants.QueryError)
	// 		}
	// 	}

	// 	watchlists, err := service.watchlistsRepository.GetUserWatchlists(spanCtx, client, users.ID, bffWatchlistsRequest.ScripId)
	// 	if err != nil {
	// 		return nil, errors.New(constants.QueryError)
	// 	}

	// 	return watchlists, nil

	// case models.DEL:

	// 	watchlistsDB, err := service.watchlistsRepository.GetValidWatchlists(spanCtx, client, users.ID, bffWatchlistsRequest.WatchlistIds)
	// 	if err != nil {
	// 		return nil, errors.New(constants.QueryError)
	// 	}

	// 	for _, wl := range watchlistsDB {

	// 		err := service.watchlistsRepository.DeleteWatchlistScrip(spanCtx, client, wl.Id, bffWatchlistsRequest.ScripId)
	// 		if err != nil {
	// 			return nil, errors.New(constants.QueryError)
	// 		}
	// 	}

	// 	watchlists, err := service.watchlistsRepository.GetUserWatchlists(spanCtx, client, users.ID, bffWatchlistsRequest.ScripId)
	// 	if err != nil {
	// 		return nil, errors.New(constants.QueryError)
	// 	}

	// 	return watchlists, nil

	case models.GET:
		var watchlists []structModels.WatchlistWithId
		if action == models.GET {
			watchlists, err = service.watchlistsRepository.GetUserWatchlists(spanCtx, client, users.ID, bffWatchlistsRequest.ScripId)
			if err != nil {
				return nil, errors.New(constants.QueryError)
			}
		}
		return watchlists, nil

	}
	return nil, err
}
