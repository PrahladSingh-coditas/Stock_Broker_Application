package business

import (
	"context"
	"errors"
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
		watchlistsWithId, err = service.watchlistRepository.GetWatchlistsWithId(ctx, tx, logger, *userId, bffAdgToWatchlistRequest.ScripId)
		if err != nil {
			return nil, nil, errors.New(constants.ErrNoWatchlistForScripMsg)
		}

		if len(watchlistsWithId) == 0 {
			return nil, nil, errors.New(constants.ErrNoWatchlistForScripMsg)
		}
	case models.DEL:
		validWatchlistIds, err := service.watchlistRepository.DeleteScripsFromWatchlists(ctx, tx, logger, *userId, bffAdgToWatchlistRequest.WatchlistIds, bffAdgToWatchlistRequest.ScripId)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}

		for _, id := range validWatchlistIds {
			watchlistsWithId = append(watchlistsWithId, models.WatchlistWithID{
				WatchlistId: id,
			})
		}

		if len(validWatchlistIds) == 0 {
			return nil, nil, errors.New(constants.ErrNoWatchlistForScripMsg)
		} else if len(validWatchlistIds) != len(bffAdgToWatchlistRequest.WatchlistIds) {
			warnings = append(warnings, "some watchlist ids were invalid")
		}

	case models.ADD:
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

		if len(result.AddedWatchlistIds)!=0 {
			for _,id:=range result.AddedWatchlistIds {
				watchlistsWithId = append(watchlistsWithId, models.WatchlistWithID{
				WatchlistId:uint64(id),
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
