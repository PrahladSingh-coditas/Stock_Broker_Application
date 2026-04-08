package business

import (
	"context"
	"errors"
	"fmt"
	"stock_broker_application/src/utils"
	"strings"
	"watchlists/commons/constants"
	"watchlists/models"
	"watchlists/repository"

	"gorm.io/gorm"
)

type AdgStoWatchlistService struct {
	adgStoWatchlistRepository repository.AdgStoWatchlistRepository
}

func NewadgStoWatchlistService(adgStoWatchlistRepository repository.AdgStoWatchlistRepository) *AdgStoWatchlistService {
	return &AdgStoWatchlistService{
		adgStoWatchlistRepository: adgStoWatchlistRepository,
	}
}

func (service *AdgStoWatchlistService) AdgStoWatchlist(ctx context.Context, username string, request models.BffAdgStoWatchlistRequest) ([]string, []models.WatchlistWithId, error) {
	postgresClient := utils.GetPostgresClient().GormDB

	ReqAction := models.Actiontype(strings.ToUpper(string(request.Action)))

	user, err := service.adgStoWatchlistRepository.GetUserByUsername(ctx, postgresClient, username)
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

		validIds, insertedIds, err := service.adgStoWatchlistRepository.AddScripToWatchlists(ctx, postgresClient, user.ID, request.ScripId, request.WatchlistIds)
		if err != nil {
			return warnings, respWatchlistWithIds, err
		}

		//no valid watchlists in the request - error
		if len(validIds) == 0 {
			return warnings, respWatchlistWithIds, constants.InvalidWatchlistsError
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

		resp, err := service.adgStoWatchlistRepository.GetWatchlistDetails(ctx, postgresClient, insertedIds)
		if err != nil {
			return warnings, respWatchlistWithIds, constants.DatabaseQueryError
		}
		for _, r := range resp {
			respWatchlistWithIds = append(respWatchlistWithIds, models.WatchlistWithId{
				Id:   r.Id,
				Name: r.WatchlistName,
			})
		}

		return warnings, respWatchlistWithIds, nil

	case models.DelAction:
		warnings := []string{}
		respWatchlistWithIds := []models.WatchlistWithId{}

		validWatchlists, err := service.adgStoWatchlistRepository.DelScripFromWatchlists(ctx, postgresClient, user.ID, request.ScripId, request.WatchlistIds)
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

		response, err := service.adgStoWatchlistRepository.GetWatchlistDetails(ctx, postgresClient, validWatchlists)
		if err != nil {
			return warnings, respWatchlistWithIds, err
		}
		for _, w := range response {
			respWatchlistWithIds = append(respWatchlistWithIds, models.WatchlistWithId{
				Id:   w.Id,
				Name: w.WatchlistName,
			})
		}
		return warnings, respWatchlistWithIds, nil

	case models.GetAction:
		warnings := []string{}
		respWatchlistWithId := []models.WatchlistWithId{}

		watchlists, err := service.adgStoWatchlistRepository.GetScripFromWatchlists(ctx, postgresClient, user.ID, request.ScripId)
		if err != nil {
			return warnings, respWatchlistWithId, err
		}
		if len(watchlists) == 0 {
			return warnings, respWatchlistWithId, constants.ScripNotInWatchlistsError
		}

		respWatchlistWithId = watchlists

		return warnings, respWatchlistWithId, nil

	default:
		return nil, nil, constants.InvalidActionTypeError
	}
}
