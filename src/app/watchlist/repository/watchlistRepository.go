package repository

import (
	"context"
	"errors"
	"stock_broker_application/src/app/watchlist/commons/constants"
	"stock_broker_application/src/app/watchlist/models"
	genericModels "stock_broker_application/src/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WatchlistRepository interface {
	GetWatchlistsWithId(ctx context.Context, db *gorm.DB, userID uint64, scripId string) ([]models.WatchlistWithID, error)
	GetUserIdByUsername(ctx context.Context, db *gorm.DB, username string) (*uint64, error)
	DeleteScripsFromWatchlists(ctx context.Context, db *gorm.DB, userID uint64, watchlistIds []uint64, scripId string) ([]uint64, error)
	AddScripsToWatchlists(ctx context.Context, db *gorm.DB, userID uint64, watchlistIds []uint64, scripId string) ([]uint64, []uint64, []uint64, error)
}

type watchlistRepository struct{}

func NewWatchlistRepository() *watchlistRepository {
	return &watchlistRepository{}
}

func (repo *watchlistRepository) GetUserIdByUsername(ctx context.Context, db *gorm.DB, username string) (*uint64, error) {
	start := time.Now()
	logger := logrus.New()

	var user genericModels.User
	result := db.WithContext(ctx).Table(constants.UsersTableName).Where(constants.UsernameCondition, username).First(&user)

	if result.RowsAffected == 0 {
		return nil, errors.New(constants.ErrUserNotFoundMsg)
	}

	if result.Error != nil {
		return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.Username: username,
		constants.Latency:  time.Since(start).Milliseconds(),
	}).Info("User ID retrieved successfully")

	return &user.ID, nil
}

func (repo *watchlistRepository) GetWatchlistsWithId(ctx context.Context, db *gorm.DB, userID uint64, scripId string) ([]models.WatchlistWithID, error) {
	var watchlistWithId []models.WatchlistWithID
	start := time.Now()
	logger := logrus.New()

	result := db.WithContext(ctx).Table(constants.WatchlistTableName).
		Select(constants.WatchlistTableName+"."+constants.FieldId+","+constants.WatchlistTableName+"."+constants.FieldwatchlistName).
		Joins("JOIN "+constants.WatchlistScripsTableName+" ON "+constants.WatchlistTableName+"."+constants.FieldId+"="+constants.WatchlistScripsTableName+"."+constants.FieldWatchlistId).
		Where(constants.WatchlistTableName+"."+constants.FielduserId+"= ? AND "+constants.WatchlistScripsTableName+"."+constants.FieldScripId+"= ?", userID, scripId).
		Scan(&watchlistWithId)

	if result.Error != nil {
		return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.UserId:  userID,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info("Watchlists with given scripId retrieved successfully")

	return watchlistWithId, nil
}

func (repo *watchlistRepository) DeleteScripsFromWatchlists(ctx context.Context, db *gorm.DB, userID uint64, watchlistIds []uint64, scripId string) ([]uint64, error) {
	var validWatchlistIds []uint64

	start := time.Now()
	logger := logrus.New()

	result := db.WithContext(ctx).
		Table(constants.WatchlistTableName).
		Where(constants.FieldId+" IN ? AND "+constants.FielduserId+" = ?", watchlistIds, userID).
		Pluck(constants.FieldId, &validWatchlistIds)

	if result.Error != nil {
		return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.UserId:  userID,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info("Watchlist IDs validated successfully")

	if len(validWatchlistIds) > 0 {

		deleteResult := db.WithContext(ctx).
			Table(constants.WatchlistScripsTableName).
			Where(constants.FieldWatchlistId+" IN ? AND "+constants.FieldScripId+" = ?", validWatchlistIds, scripId).
			Delete(nil)

		if deleteResult.Error != nil {
			return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
		}

		logger.WithFields(logrus.Fields{
			constants.UserId:  userID,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info("Scrip removed from watchlists successfully")

		if deleteResult.RowsAffected > 0 {
			result = db.WithContext(ctx).
				Table(constants.WatchlistTableName).
				Where(constants.FieldId+" IN ?", validWatchlistIds).
				Updates(map[string]interface{}{
					constants.FieldScripCount:  gorm.Expr("scrip_count - 1"),
					constants.FieldLastUpdated: gorm.Expr("NOW()"),
				})

			if result.Error != nil {
				return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
			}

			logger.WithFields(logrus.Fields{
				constants.UserId:  userID,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info("Watchlist scrip count updated successfully")
		}
	}

	return validWatchlistIds, nil
}

func (repo *watchlistRepository) AddScripsToWatchlists(ctx context.Context, db *gorm.DB, userID uint64, watchlistIds []uint64, scripId string) ([]uint64, []uint64, []uint64, error) {
	var validWatchlistIds []uint64
	var count int64
	start := time.Now()
	logger := logrus.New()

	//Checking the scrip ID here.
	result := db.WithContext(ctx).
		Table("scrip_masters").
		Where("id = ?", scripId).
		Count(&count)

	if result.Error != nil {
		return nil, nil, nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	if count == 0 {
		return nil, nil, nil, errors.New(constants.ErrScripNotFoundMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.UserId:  userID,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info("Scrip ID validated successfully")

	//Checking the watchlist ids and ownership here.
	result = db.WithContext(ctx).
		Table(constants.WatchlistTableName).
		Where(constants.FieldId+" IN ? AND "+constants.FielduserId+" = ?", watchlistIds, userID).
		Pluck(constants.FieldId, &validWatchlistIds)

	if result.Error != nil {
		return nil, nil, nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	if len(validWatchlistIds) == 0 {
		return nil, nil, nil, errors.New(constants.ErrNoWatchlistForScripMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.UserId:  userID,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info("Watchlist IDs validated successfully")

	//Adding the scrip to valid watchlists here.
	var existingWatchlistIds []uint64

	result = db.WithContext(ctx).
		Table(constants.WatchlistScripsTableName).
		Where(constants.FieldWatchlistId+" IN ? AND "+constants.FieldScripId+" = ?", validWatchlistIds, scripId).
		Pluck(constants.FieldWatchlistId, &existingWatchlistIds)

	if result.Error != nil {
		return nil, nil, nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	existingWatchlistIdsMap := make(map[uint64]bool)

	for _, id := range existingWatchlistIds {
		existingWatchlistIdsMap[id] = true
	}

	type Watchlist struct {
		WatchlistId uint64 `gorm:"column:id"`
		ScripCount  int    `gorm:"column:scrip_count"`
	}

	var watchlists []Watchlist

	result = db.WithContext(ctx).Table(constants.WatchlistTableName).
		Select(constants.WatchlistTableName+"."+constants.FieldId+", "+constants.WatchlistTableName+"."+constants.FieldScripCount).
		Where(constants.WatchlistTableName+"."+constants.FieldId+" IN ?", validWatchlistIds).
		Scan(&watchlists)

	var addedWatchlistIds []uint64
	var SkipppedWatchlistIds []uint64
	var limitExceededWatchlistIds []uint64

	for _, w := range watchlists {

		if existingWatchlistIdsMap[w.WatchlistId] {
			SkipppedWatchlistIds = append(SkipppedWatchlistIds, w.WatchlistId)
			continue
		}

		if w.ScripCount >= constants.MaxScripsPerWatchlist {
			limitExceededWatchlistIds = append(limitExceededWatchlistIds, w.WatchlistId)
			continue
		}

		result := db.WithContext(ctx).
			Table(constants.WatchlistScripsTableName).
			Create(map[string]interface{}{
				constants.FieldWatchlistId: w.WatchlistId,
				constants.FieldScripId:     scripId,
			})

		if result.Error != nil {
			return nil, nil, nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
		}

		result = db.WithContext(ctx).
			Table(constants.WatchlistTableName).
			Where(constants.FieldId+" = ?", w.WatchlistId).
			Update(constants.FieldScripCount, gorm.Expr("scrip_count + 1"))

		if result.Error != nil {
			return nil, nil, nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
		}
		addedWatchlistIds = append(addedWatchlistIds, w.WatchlistId)
		logger.WithFields(logrus.Fields{
			constants.UserId:  userID,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info("Scrip added to watchlists successfully")
	}

	return addedWatchlistIds, SkipppedWatchlistIds, limitExceededWatchlistIds, nil
}
