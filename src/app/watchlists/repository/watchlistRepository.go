package repository

import (
	"context"
	"errors"
	"fmt"
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
	GetValidWatchlists(ctx context.Context, db *gorm.DB, userId uint64, watchlistIds []uint64) ([]genericModels.Watchlists, error)
	CheckDuplicate(ctx context.Context, db *gorm.DB, watchlistId uint64, scripId string) (bool, error)
	InsertWatchlistScrip(ctx context.Context, db *gorm.DB, watchlistId uint64, scripId string) error
	DeleteWatchlistScrip(ctx context.Context, db *gorm.DB, watchlistId uint64, scripId string) error
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

func (user *watchlistsRepository) GetUserWatchlists(ctx context.Context, db *gorm.DB, userId uint64, scripId string) ([]structModels.WatchlistWithId, error) {

	start := time.Now()
	logger := logrus.New()

	var result []structModels.WatchlistWithId

	err := db.Table(constants.WatclistsTableName+" as w").
		Select("w."+constants.FieldWatchlistId+" as watchlist_id, w."+constants.FieldWatchlistName+" as watchlist_name").
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

	var scripmasterDB genericModels.ScripMaster

	result := db.WithContext(ctx).
		Table(constants.ScripMasterTableName).
		Where(constants.FieldScripId, scripId).
		First(&scripmasterDB)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New(constants.ScripIdNotFoundError)
		}
		return false, fmt.Errorf("%s: %w", constants.AuthenticationFailedError, result.Error)
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.ScripIdExistsSuccessMsg)

	return true, nil

}

func (repo *watchlistsRepository) GetValidWatchlists(ctx context.Context, db *gorm.DB, userId uint64, watchlistIds []uint64) ([]genericModels.Watchlists, error) {

	var watchlists []genericModels.Watchlists

	err := db.Table(constants.WatclistsTableName).
		Where(constants.FieldWatchUserId+" = ? AND "+constants.FieldWatchlistId+" IN ?", userId, watchlistIds).
		Find(&watchlists).Error

	return watchlists, err
}
func (repo *watchlistsRepository) CheckDuplicate(ctx context.Context, db *gorm.DB, watchlistId uint64, scripId string) (bool, error) {

	var count int64

	err := db.Table(constants.WatchScripTableName).
		Where(constants.FieldWatchId+" = ? AND "+constants.FieldWScripId+" = ?", watchlistId, scripId).
		Count(&count).Error

	return count > 0, err
}
func (repo *watchlistsRepository) InsertWatchlistScrip(ctx context.Context, db *gorm.DB, watchlistId uint64, scripId string) error {

	insert := genericModels.WatchlistScrip{
		WatchlistId: watchlistId,
		ScripId:     scripId,
	}

	return db.Table(constants.WatchScripTableName).Create(&insert).Error
}
func (repo *watchlistsRepository) DeleteWatchlistScrip(ctx context.Context, db *gorm.DB, watchlistId uint64, scripId string) error {

	return db.Table(constants.WatchScripTableName).
		Where(constants.FieldWatchId+" = ? AND "+constants.FieldWScripId+" = ?", watchlistId, scripId).
		Delete(&genericModels.WatchlistScrip{}).Error
}
