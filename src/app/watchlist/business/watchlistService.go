package business

import (
	"context"
	"errors"
	"fmt"
	"stock_broker_application/src/app/watchlist/commons/constants"
	"stock_broker_application/src/app/watchlist/models"
	"stock_broker_application/src/app/watchlist/repository"
	"stock_broker_application/src/utils"

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
	postgresClinet := utils.GetPostgresClient()
	tx := postgresClinet.GormDB.Begin()
	var watchlistsWithId []models.WatchlistWithID
	var warnings []string

	userId, err := service.watchlistRepository.GetUserIdByUsername(ctx, tx, logger, username)

	if err != nil {
		return nil, nil, err
	}

	switch bffAdgToWatchlistRequest.Action {
	case models.GET:
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
	case models.DEL:

		if len(bffAdgToWatchlistRequest.WatchlistIds) == 0 {
			return nil, nil, errors.New("watchlist Ids can't be empty for this operation")
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

		fmt.Println(resultListsForDEL.ValidWatchlistIds)

		if len(resultListsForDEL.ValidWatchlistIds) == 0 {
			return nil, nil, errors.New(constants.ErrNoWatchlistForScripMsg)
		} else if len(resultListsForDEL.ValidWatchlistIds) != len(bffAdgToWatchlistRequest.WatchlistIds) {
			warnings = append(warnings, "some watchlist ids were invalid")
		}

	case models.ADD:

		if len(bffAdgToWatchlistRequest.WatchlistIds) == 0 {
			return nil, nil, errors.New("watchlist Ids can't be empty for this operation")
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
	}

	if len(warnings) == 0 {
		warnings = append(warnings, constants.NoWarningsMsg)
	}

	tx.Commit()

	return watchlistsWithId, warnings, nil
}
