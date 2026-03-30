package repository

import (
	"context"
	"errors"
	
	"time"
	"watchlist/commons/constants"
	"watchlist/models"

	genericModels "stock_broker_application/src/models"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WatchlistRepository interface {
	GetUserFromDb(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error)
	WatchlistGetOperation(ctx context.Context, db *gorm.DB, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, error)
	WatchlistAddOperation(ctx context.Context, db *gorm.DB, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, []models.WatchlistValidation, error)
	WatchlistDeleteOperation(ctx context.Context, db *gorm.DB, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, []models.WatchlistValidation, error)
}

type watchlistRepository struct{}

func NewWatchlistRepository() *watchlistRepository {
	return &watchlistRepository{}
}

func (user *watchlistRepository) GetUserFromDb(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error) {
	start := time.Now()
	logger := logrus.New()

	var ExistingUser genericModels.User

	result := db.WithContext(ctx).
		Table("users").
		Where(constants.UsernameField, username).
		First(&ExistingUser)

	if result.Error != nil {
		return nil, result.Error
	}

	logger.WithFields(logrus.Fields{
		"user":    username,
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserLoggedInSuccessMsg)

	return &ExistingUser, nil

}

func (user *watchlistRepository) WatchlistGetOperation(ctx context.Context, db *gorm.DB, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, error) {
	start := time.Now()
	logger := logrus.New()

	var WatchlistNamewithId []models.WatchlistWithId

	result := db.WithContext(ctx).
		Table("watchlists w").
		Select(" DISTINCT w.id AS watchlist_id, w.watchlist_name as watchlist_name").
		Joins("JOIN watchlist_scrips ws ON w.id = ws.watchlist_id").
		Where("w.user_id= ? and ws.scrip_id=?", UserId, bffAdgToWatchlistRequest.ScripId).
		Scan(&WatchlistNamewithId)

	if result.Error != nil {
		return nil, result.Error
	}

	logger.WithFields(logrus.Fields{
		"user":    UserId,
		"latency": time.Since(start).Milliseconds(),
	}).Info("Watchlists fetched successfully")

	return WatchlistNamewithId, nil

}
func (repo *watchlistRepository) ValidateWatchlistOperation(ctx context.Context, db *gorm.DB, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistValidation, error) {

	var results []models.WatchlistValidation

	err := db.Raw(
		constants.WatchlistValidateQuery,
		pq.Array(bffAdgToWatchlistRequest.WatchlistIds), // UNNEST
		bffAdgToWatchlistRequest.ScripId,                // scrip_exists
		bffAdgToWatchlistRequest.ScripId,                // scrip_present_check
		UserId,                                          // CanInsert ownership
		UserId,                                          // CanDelete ownership

		//reasons
		constants.ErrWatchlistNotFound,
		UserId,
		constants.ErrWatchlistsNotOfUser,
		constants.ErrScripAlreadyExists,
		constants.ErrScripNotPresent,
		constants.ErrScripLimitExceeded,
	).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *watchlistRepository) WatchlistAddOperation(ctx context.Context, db *gorm.DB, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, []models.WatchlistValidation, error) {

	start := time.Now()
	logger := logrus.New()

	validationResults, err := repo.ValidateWatchlistOperation(ctx, db, UserId, bffAdgToWatchlistRequest)
	if err != nil {
		return nil, nil, errors.New(constants.ErrValidatingConstraints)
	}

	var validWatchlists []uint64
	var warnings []models.WatchlistValidation

	for _, v := range validationResults {
		if v.CanInsert {
			validWatchlists = append(validWatchlists, v.WatchlistId)
		} else {
			warnings = append(warnings, v)
		}
	}

	if len(validWatchlists) == 0 {
		return []models.WatchlistWithId{}, warnings, nil
	}

	// append valid watchlists and their scipr ids to slice of struct
	var records []models.WatchlistAndScrip
	for _, id := range validWatchlists {
		records = append(records, models.WatchlistAndScrip{
			WatchlistId: id,
			ScripId:     bffAdgToWatchlistRequest.ScripId,
		})
	}

	var result []models.WatchlistWithId
	err = db.WithContext(ctx).
		Raw(constants.WatchlistAddQuery, bffAdgToWatchlistRequest.ScripId, validWatchlists).Scan(&result).Error

	if err != nil {
		return nil, nil, err
	}

	logger.WithFields(logrus.Fields{
		"scripId": bffAdgToWatchlistRequest.ScripId,
		"latency": time.Since(start).Milliseconds(),
	}).Info("Scrip added successfully")

	return result, warnings, nil
}

func (repo *watchlistRepository) WatchlistDeleteOperation(ctx context.Context, db *gorm.DB, UserId uint64, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) ([]models.WatchlistWithId, []models.WatchlistValidation, error) {
	start := time.Now()
	logger := logrus.New()

	validations, err := repo.ValidateWatchlistOperation(ctx, db, UserId, bffAdgToWatchlistRequest)
	if err != nil {
		return nil, nil, errors.New(constants.ErrValidatingConstraints)
	}

	var validWatchlistIds []uint64
	var warnings []models.WatchlistValidation

	for _, val := range validations {
		if val.CanDelete {
			validWatchlistIds = append(validWatchlistIds, val.WatchlistId)
		} else {
			warnings = append(warnings, val)
		}
	}
	if len(validWatchlistIds) == 0 {
		return []models.WatchlistWithId{}, warnings, nil
	}

	
	var deletedWatchlists []models.WatchlistWithId

	err = db.WithContext(ctx).
		Raw(constants.WatchlistDeleteQuery, validWatchlistIds, bffAdgToWatchlistRequest.ScripId).Scan(&deletedWatchlists).Error
	if err != nil {
		return nil, nil, err
	}

	logger.WithFields(logrus.Fields{
		"scripId": bffAdgToWatchlistRequest.ScripId,
		"latency": time.Since(start).Milliseconds(),
	}).Info("Scrip deleted successfully")

	return deletedWatchlists, warnings, nil

}
