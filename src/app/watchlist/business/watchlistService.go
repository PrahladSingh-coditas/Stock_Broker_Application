package business

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

	actionType := bffAdgToWatchlistRequest.Action

	//get user from db
	user, err := service.watchlistRepository.GetUserFromDb(spanCtx, client, username)
	if err != nil {
		return nil, nil, errors.New("Couldnt find user")
	}

	switch actionType {
	case models.ADD:

		if strings.TrimSpace(bffAdgToWatchlistRequest.ScripId) == "" {
			return nil, nil, errors.New(constants.ErrEmptyScripId)
		}

		if len(bffAdgToWatchlistRequest.WatchlistIds) == 0 {
			return nil, nil, errors.New(constants.ErrEmptyWatchlists)
		}

		addedWatchlists, warningsResult, err := service.watchlistRepository.WatchlistAddOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest.ScripId, bffAdgToWatchlistRequest.WatchlistIds)
		if err != nil {
			return nil, nil, errors.New(constants.ErrNoRowsAffected)
		}

		var warnings []string
		for _, w := range warningsResult {
			msg := "Watchlist " + strconv.FormatUint(w.WatchlistId, 10) + ": " + w.Reason
			warnings = append(warnings, msg)
		}
		if len(addedWatchlists) == 0 {
			return addedWatchlists, warnings, nil
		}

		return addedWatchlists, warnings, nil

	case models.DEL:
		if strings.TrimSpace(bffAdgToWatchlistRequest.ScripId) == "" {
			return nil, nil, errors.New(constants.ErrEmptyScripId)
		}

		if len(bffAdgToWatchlistRequest.WatchlistIds) == 0 {
			return nil, nil, errors.New(constants.ErrEmptyWatchlists)
		}

		deletedWatchlists, warningsResult, err := service.watchlistRepository.WatchlistDeleteOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest.ScripId, bffAdgToWatchlistRequest.WatchlistIds)
		if err != nil {
			return nil, nil, err
		}

		var warnings []string
		for _, w := range warningsResult {
			msg := "Watchlist " + strconv.FormatUint(w.WatchlistId, 10) + ": " + w.Reason
			warnings = append(warnings, msg)
		}

		fmt.Println("warnings:", warnings)
		fmt.Println("deleted:", deletedWatchlists)
		if len(deletedWatchlists) == 0 {
			return deletedWatchlists, nil, nil
		}
		

		return deletedWatchlists, warnings, nil

	case models.GET:

		if strings.TrimSpace(bffAdgToWatchlistRequest.ScripId) == "" {
			return nil, nil, errors.New(constants.ErrEmptyScripId)
		}

		WatchlistNamewithId, err := service.watchlistRepository.WatchlistGetOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest.ScripId)
		if err != nil {
			return nil, nil, errors.New(constants.ErrNoRowsAffected)
		}

		if len(WatchlistNamewithId) == 0 {
			return nil, nil, errors.New(constants.ErrWatchlistNotFound)
		}

		return WatchlistNamewithId, nil, nil

	default:
		return nil, nil, errors.New(constants.ErrInvalidAction)

	}

	return nil, nil, err
}
