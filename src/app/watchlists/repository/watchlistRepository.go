package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	genericModels "stock_broker_application/src/models"
	"time"
	"watchlists/commons/constants"
	structModels "watchlists/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WatchlistsRepository interface {
	GetUserId(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error)
	GetUserWatchlists(ctx context.Context, db *gorm.DB, userId uint64, scripId string) ([]structModels.WatchlistWithId, error)
	CheckScripExists(ctx context.Context, db *gorm.DB, scripId string) (bool, error)
	GetValidWatchlists(ctx context.Context, db *gorm.DB, userId uint64, watchlistIds []uint64) ([]genericModels.Watchlists, []uint64, error)
	CheckDuplicate(ctx context.Context, db *gorm.DB, capNotFullIds []uint64, scripId string) ([]uint64, []uint64, error)
	InsertWatchlistScrip(ctx context.Context, db *gorm.DB, notDuplicates []uint64, scripId string) ([]uint64, error)
	GetWatchlistDetails(ctx context.Context, db *gorm.DB, addedIds []uint64) ([]genericModels.Watchlists, error)
	DeleteWatchlistScrip(ctx context.Context, db *gorm.DB, watchlistIdNames []structModels.WatchlistWithId, scripId string) error
}

type watchlistsRepository struct{}

func NewWatchlistsRepository() *watchlistsRepository {
	return &watchlistsRepository{}
}

func (user *watchlistsRepository) GetUserId(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error) {

	start := time.Now()
	logger := logrus.New()

	var userDB genericModels.User

	result := db.WithContext(ctx).
		Table(constants.UsersTableName).
		Where(constants.FieldUsername, username).
		First(&userDB)

	//we have checked here with username
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.UserNotFoundError)
		}
		return nil, fmt.Errorf("%s: %w", constants.AuthenticationFailedError, result.Error)
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.GetUserIdSuccessMsg)

	return &userDB, nil
}

func (repo *watchlistsRepository) GetUserWatchlists(ctx context.Context, db *gorm.DB, userId uint64, scripId string) ([]structModels.WatchlistWithId, error) {

	start := time.Now()
	logger := logrus.New()

	var result []structModels.WatchlistWithId

	err := db.WithContext(ctx).
		Table(constants.WatclistsTableName+" as w").
		Select("DISTINCT w."+constants.FieldWatchlistId+" as watchlist_id, w."+constants.FieldWatchlistName+" as watchlist_name").
		Joins("JOIN "+constants.WatchScripTableName+" as ws ON ws."+constants.FieldWatchId+" = w."+constants.FieldWatchlistId).
		Where("w."+constants.FieldWatchUserId+" = ? AND ws."+constants.FieldWScripId+" = ?", userId, scripId).
		Scan(&result).Error

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.GetWatchlistIdsSuccessMsg)

	return result, err
}

func (user *watchlistsRepository) CheckScripExists(ctx context.Context, db *gorm.DB, scripId string) (bool, error) {
	start := time.Now()
	logger := logrus.New()

	var count int64
	fmt.Println("Checking existence for scripId:", scripId, "hhhj")
	err := db.WithContext(ctx).
		Table(constants.ScripMasterTableName).
		Where("id = ?", scripId).
		Count(&count).Error

	if err != nil {
		log.Println("Error:", err)
	}

	fmt.Println("Count:", count)

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
		"count":   count,
	}).Info(constants.ScripIdExistsSuccessMsg)

	return count > 0, nil
}

func (repo *watchlistsRepository) GetValidWatchlists(ctx context.Context, db *gorm.DB, userId uint64, watchlistIds []uint64) ([]genericModels.Watchlists, []uint64, error) {

	var watchlists []genericModels.Watchlists

	err := db.WithContext(ctx).
		Table(constants.WatclistsTableName).
		Where("user_id = ? AND id IN ?", userId, watchlistIds).
		Find(&watchlists).Error

	if err != nil {
		return nil, nil, err
	}

	validMap := make(map[uint64]struct{})
	for _, wl := range watchlists {
		validMap[wl.Id] = struct{}{}
	}

	var notValid []uint64
	for _, id := range watchlistIds {
		if _, ok := validMap[id]; !ok {
			notValid = append(notValid, id)
		}
	}

	return watchlists, notValid, nil
}

func (repo *watchlistsRepository) CheckDuplicate(ctx context.Context, db *gorm.DB, ids []uint64, scripId string) ([]uint64, []uint64, error) {

	var existing []genericModels.WatchlistScrips

	err := db.WithContext(ctx).
		Table(constants.WatchScripTableName).
		Where("watchlist_id IN ? AND scrip_id = ?", ids, scripId).
		Find(&existing).Error

	if err != nil {
		return nil, nil, err
	}

	dupMap := make(map[uint64]struct{})
	for _, e := range existing {
		dupMap[e.WatchlistId] = struct{}{}
	}

	var duplicates, notDuplicates []uint64
	for _, id := range ids {
		if _, ok := dupMap[id]; ok {
			duplicates = append(duplicates, id)
		} else {
			notDuplicates = append(notDuplicates, id)
		}
	}

	return duplicates, notDuplicates, nil
}

func (repo *watchlistsRepository) InsertWatchlistScrip(ctx context.Context, db *gorm.DB, ids []uint64, scripId string) ([]uint64, error) {

	var added []uint64

	for _, id := range ids {
		var wl struct {
			ScripCount int
		}
		err := db.WithContext(ctx).
			Table(constants.WatclistsTableName).
			Where("id = ?", id).
			Select("scrip_count").
			Take(&wl).Error
		if err != nil {
			return nil, err
		}

		if wl.ScripCount >= 10 {
			continue
		}

		err = db.WithContext(ctx).
			Table(constants.WatchScripTableName).
			Create(&genericModels.WatchlistScrips{
				WatchlistId: id,
				ScripId:     scripId,
			}).Error
		if err != nil {
			return nil, err
		}

		err = db.WithContext(ctx).
			Table(constants.WatclistsTableName).
			Where("id = ?", id).
			Update("scrip_count", gorm.Expr("scrip_count + 1")).Error
		if err != nil {
			return nil, err
		}

		added = append(added, id)
	}

	return added, nil
}

func (repo *watchlistsRepository) GetWatchlistDetails(ctx context.Context, db *gorm.DB, ids []uint64) ([]genericModels.Watchlists, error) {

	var watcheslists []genericModels.Watchlists

	err := db.WithContext(ctx).
		Table(constants.WatclistsTableName).
		Where("id IN ?", ids).
		Find(&watcheslists).Error

	return watcheslists, err
}

func (repo *watchlistsRepository) DeleteWatchlistScrip(ctx context.Context, db *gorm.DB, watchlistIdNames []structModels.WatchlistWithId, scripId string) error {
	var watchlistIds []uint64
	for _, wl := range watchlistIdNames {
		watchlistIds = append(watchlistIds, uint64(wl.Watchlist_ID))
	}

	err := db.WithContext(ctx).
		Table(constants.WatchScripTableName).
		Where(constants.FieldWatchId+" IN ? AND "+constants.FieldWScripId+" = ?", watchlistIds, scripId).
		Delete(&genericModels.WatchlistScrips{}).Error

	if err != nil {
		return err
	}
	err = db.WithContext(ctx).
		Table(constants.WatclistsTableName).
		Where("id IN ?", watchlistIds).
		Updates(map[string]interface{}{
			"scrip_count":       gorm.Expr("scrip_count - 1"),
			"last_updated  _at": gorm.Expr("NOW()"),
		}).Error
	if err != nil {
		return err
	}
	return nil
}
