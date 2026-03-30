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

	//get user from db and to extract and use user.ID in ADG operations
	user, err := service.watchlistRepository.GetUserFromDb(spanCtx, client, username)
	if err != nil {
		return nil, nil, errors.New(constants.ErrUserNotFound)
	}

	switch actionType {
	case models.ADD:

		// if strings.TrimSpace(bffAdgToWatchlistRequest.ScripId) == "" {
		// 	return nil, nil, errors.New(constants.ErrEmptyScripId)
		// }

		// if len(bffAdgToWatchlistRequest.WatchlistIds) == 0 {
		// 	return nil, nil, errors.New(constants.ErrEmptyWatchlists)
		// }

		addedWatchlists, warningsResult, err := service.watchlistRepository.WatchlistAddOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
		if err != nil {
			return nil, nil, err
		}

		var warningsAll []string
		for _, w := range warningsResult {
			msg := "Watchlist " + strconv.FormatUint(w.WatchlistId, 10) + ": " + w.Reason
			warningsAll = append(warningsAll, msg)
		}
		
		var warningsNotOfUser []string
		for _,val := range warningsAll{
			if strings.Contains(val,constants.ErrWatchlistsNotOfUser){
				warningsNotOfUser = append(warningsNotOfUser, val)
			} else {
				continue
			}
		}

		if len(warningsNotOfUser) == len(bffAdgToWatchlistRequest.WatchlistIds){
			return nil,nil,errors.New(constants.ErrWatchlistsNotOfUser)
		} else{
			return addedWatchlists,warningsAll,nil
		}

	case models.DEL:
		// if strings.TrimSpace(bffAdgToWatchlistRequest.ScripId) == "" {
		// 	return nil, nil, errors.New(constants.ErrEmptyScripId)
		// }

		// if len(bffAdgToWatchlistRequest.WatchlistIds) == 0 {
		// 	return nil, nil, errors.New(constants.ErrEmptyWatchlists)
		// }

		deletedWatchlists, warningsResult, err := service.watchlistRepository.WatchlistDeleteOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
		if err != nil {
			return nil, nil, err
		}

		var warningsAll []string
		for _, w := range warningsResult {
			msg := "Watchlist " + strconv.FormatUint(w.WatchlistId, 10) + ": " + w.Reason
			warningsAll = append(warningsAll, msg)
		}

		fmt.Println("warnings:", warningsAll)
		fmt.Println("deleted:", deletedWatchlists)
		

		var warningsNotOfUser []string
		for _,val := range warningsAll{
			if strings.Contains(val,constants.ErrWatchlistsNotOfUser){
				warningsNotOfUser = append(warningsNotOfUser, val)
			} else {
				continue
			}
		}
		
		if len(warningsNotOfUser) == len(bffAdgToWatchlistRequest.WatchlistIds){
			return nil,nil,errors.New(constants.ErrWatchlistsNotOfUser)
		} else{
			return deletedWatchlists,warningsAll,nil
		}


	case models.GET:

		// if strings.TrimSpace(bffAdgToWatchlistRequest.ScripId) == "" {
		// 	return nil, nil, errors.New(constants.ErrEmptyScripId)
		// }

		WatchlistNamewithId, err := service.watchlistRepository.WatchlistGetOperation(spanCtx, client, user.ID, bffAdgToWatchlistRequest)
		if err != nil {
			return nil, nil, errors.New(constants.ErrNoRowsAffected)
		}

		if len(WatchlistNamewithId) == 0 {
			return nil, nil, errors.New(constants.ErrWatchlistNotFound)
		}

		var warnings []string
		warnings = append(warnings,constants.ErrWatchlistNotRequired)
		if len(bffAdgToWatchlistRequest.WatchlistIds) != 0 {
			return WatchlistNamewithId, warnings,nil
		}

		return WatchlistNamewithId, nil, nil

	default:
		return nil, nil, errors.New(constants.ErrInvalidAction)

	}

}
