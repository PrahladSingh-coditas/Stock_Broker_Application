package buisness

import (
	"errors"
	"watchList/models"
	"watchList/repository"
)

type WatchlistService interface {
	HandleWatchlist(username string, req models.BFFAdgToWatchListRequest) (*models.BFFAdgToWatchlistResponse, error)
}

type watchlistService struct {
	repo repository.WatchListRepository
}

func NewWatchlistService(repo repository.WatchListRepository) *watchlistService {
	return &watchlistService{repo: repo}
}

func (controller *watchlistService) HandleWatchlist(username string, req models.BFFAdgToWatchListRequest) (*models.BFFAdgToWatchlistResponse, error) {

	user, err := controller.repo.GetUserByUsername(username)

	if err != nil {
		return nil, errors.New("user not found")
	}

	switch req.Action {

	case models.ADD:
		return controller.handleAdd(user.ID, req)

	default:
		return nil, errors.New("invalid action")
	}
}

func (controller *watchlistService) handleAdd(userID uint64, req models.BFFAdgToWatchListRequest) (*models.BFFAdgToWatchlistResponse, error) {
	warnings := make([]string, 0)
	result := make([]models.WatchListWithId, 0)

	//fetching watchlist belong to user
	watchlists, err := controller.repo.GetWatchlistsByIDs(userID, req.WatchlistIds)
	if err != nil {
		return nil, err
	}

	for _, wl := range watchlists {

		//checking duplicate
		alreadyexists, err := controller.repo.IsScripInWatchlist(wl.ID, req.ScripId)
		if err != nil {
			return nil, err
		}

		if alreadyexists {
			warnings = append(warnings, "scrip already exists in "+wl.WatchlistName)
			continue
		}

		//adding scrip to watchlist
		err = controller.repo.AddScripToWatchlist(wl.ID, req.ScripId)
		if err != nil {
			if err.Error() == "watchlist limit reached" {
				warnings = append(warnings,
					"watchlist "+wl.WatchlistName+" reached max limit")
				continue
			}
			return nil, err
		}

		result = append(result, models.WatchListWithId{
			ID:            wl.ID,
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
