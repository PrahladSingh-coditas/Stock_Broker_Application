package buisness

import (
	"context"
	"errors"
	"watchList/models"
	"watchList/repository"
)

type WatchlistService interface {
	HandleWatchlist(ctx context.Context, username string, req models.BFFAdgToWatchListRequest) (*models.BFFAdgToWatchlistResponse, error)
}

type watchlistService struct {
	repo repository.WatchListRepository
}

func NewWatchlistService(repo repository.WatchListRepository) *watchlistService {
	return &watchlistService{repo: repo}
}

func (controller *watchlistService) HandleWatchlist(ctx context.Context, username string, req models.BFFAdgToWatchListRequest) (*models.BFFAdgToWatchlistResponse, error) {

	user, err := controller.repo.GetUserByUsername(ctx, username)

	if err != nil {
		return nil, errors.New("user not found")
	}

	switch req.Action {

	case models.ADD:
		return controller.handleAdd(ctx, user.ID, req)
	case models.DEL:
		return controller.handleDelete(ctx, user.ID, req)
	case models.GET:
		return controller.handleGet(ctx, user.ID, req)

	default:
		return nil, errors.New("invalid action")
	}
}

func (controller *watchlistService) handleAdd(ctx context.Context, userID uint64, req models.BFFAdgToWatchListRequest) (*models.BFFAdgToWatchlistResponse, error) {
	warnings := make([]string, 0)
	result := make([]models.WatchListWithId, 0)

	//fetching watchlist belong to user
	watchlists, err := controller.repo.GetWatchlistsByIDs(ctx, userID, req.WatchlistIds)
	if err != nil {
		return nil, err
	}

	for _, wl := range watchlists {

		//checking duplicate
		alreadyexists, err := controller.repo.IsScripInWatchlist(ctx, wl.Id, req.ScripId)
		if err != nil {
			return nil, err
		}

		if alreadyexists {
			warnings = append(warnings, "scrip already exists in "+wl.WatchlistName)
			continue
		}

		//adding scrip to watchlist
		err = controller.repo.AddScripToWatchlist(ctx, wl.Id, req.ScripId)
		if err != nil {
			if err.Error() == "watchlist limit reached" {
				warnings = append(warnings,
					"watchlist "+wl.WatchlistName+" reached max limit")
				continue
			}
			return nil, err
		}

		result = append(result, models.WatchListWithId{
			ID:            wl.Id,
			WatchListName: wl.WatchlistName,
		})

	}

	return &models.BFFAdgToWatchlistResponse{
		Status:          "success",
		Action:          models.ADD,
		WatchlistWithId: result,
		Warning:         warnings,
	}, nil
}

func (controller *watchlistService) handleDelete(ctx context.Context, userID uint64, req models.BFFAdgToWatchListRequest) (*models.BFFAdgToWatchlistResponse, error) {

	result := make([]models.WatchListWithId, 0)
	warnings := make([]string, 0)

	// Validate watchlists belong to user
	watchlists, err := controller.repo.GetWatchlistsByIDs(ctx, userID, req.WatchlistIds)
	if err != nil {
		return nil, err
	}

	if len(watchlists) != len(req.WatchlistIds) {
		return nil, errors.New("one or more watchlists do not belong to user")
	}

	// Loop through watchlists
	for _, wl := range watchlists {

		// Check if scrip exists in watchlist
		exists, err := controller.repo.IsScripInWatchlist(ctx, wl.Id, req.ScripId)
		if err != nil {
			return nil, err
		}

		if !exists {
			warnings = append(warnings,
				"scrip not found in "+wl.WatchlistName)
			continue
		}

		// Delete + decrement count
		err = controller.repo.DeleteScripFromWatchlist(ctx, wl.Id, req.ScripId)
		if err != nil {
			return nil, err
		}

		result = append(result, models.WatchListWithId{
			ID:            wl.Id,
			WatchListName: wl.WatchlistName,
		})
	}

	return &models.BFFAdgToWatchlistResponse{
		Status:          "success",
		Action:          models.DEL,
		WatchlistWithId: result,
		Warning:         warnings,
	}, nil

}

func (controller *watchlistService) handleGet(ctx context.Context, userID uint64, req models.BFFAdgToWatchListRequest) (*models.BFFAdgToWatchlistResponse, error) {

	result := make([]models.WatchListWithId, 0)

	watchlists, err := controller.repo.GetWatchlistsContainingScrip(ctx, userID, req.ScripId)
	if err != nil {
		return nil, err
	}

	for wl := range watchlists {
		result = append(result, models.WatchListWithId{
			ID:            watchlists[wl].Id,
			WatchListName: watchlists[wl].WatchlistName,
		})
	}

	return &models.BFFAdgToWatchlistResponse{
		Status:          "success",
		Action:          models.GET,
		WatchlistWithId: result,
		Warning:         []string{},
	}, nil
}
