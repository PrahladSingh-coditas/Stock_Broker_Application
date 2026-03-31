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

	"github.com/lib/pq"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WatchlistsRepository interface {
	GetUserId(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error)
	GetUserWatchlists(ctx context.Context, db *gorm.DB, userId uint64, scripId string) ([]structModels.WatchlistWithId, error)
	CheckScripExists(ctx context.Context, db *gorm.DB, scripId string) (bool, error)
	GetValidWatchlists(ctx context.Context, db *gorm.DB, userId uint64, watchlistIds []uint64) ([]genericModels.Watchlists, []uint64, error)

	AddScripWithCTE(ctx context.Context, db *gorm.DB, userId uint64, watchlistIds []uint64, scripId string) ([]structModels.WatchlistWithId, error)
	DeleteScripWithCTE(ctx context.Context, db *gorm.DB, userId uint64, watchlistIds []uint64, scripId string) ([]structModels.WatchlistWithId, error)
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
	err := db.WithContext(ctx).
		Table(constants.ScripMasterTableName).
		Where("id = ?", scripId).
		Count(&count).Error

	if err != nil {
		log.Println("Error:", err)
	}

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

	fmt.Println("Fetched watchlists from DB:", watchlists)

	validMap := make(map[uint64]bool)
	for _, wl := range watchlists {
		validMap[wl.Id] = true
	}

	fmt.Println("Valid watchlist IDs:", validMap)
	var notValid []uint64
	for _, id := range watchlistIds {
		if _, ok := validMap[id]; !ok {
			notValid = append(notValid, id)
		}
	}

	return watchlists, notValid, nil
}

func (repo *watchlistsRepository) AddScripWithCTE(ctx context.Context, db *gorm.DB, userId uint64, watchlistIds []uint64, scripId string) ([]structModels.WatchlistWithId, error) {

	var result []structModels.WatchlistWithId

	query := `
    WITH inserted AS (
		INSERT INTO watchlist_scrips (watchlist_id, scrip_id)
		SELECT wl.id, ?
		FROM watchlists wl
		WHERE wl.user_id = ?
		AND wl.id = ANY(?)
		RETURNING watchlist_id
	),
	updated AS (
		UPDATE watchlists
		SET scrip_count = scrip_count + 1,
			last_updated_at = NOW()
		WHERE id IN (SELECT watchlist_id FROM inserted)
		RETURNING id, watchlist_name
	)
	SELECT id AS watchlist_id, watchlist_name
	FROM updated;
    `

	err := db.WithContext(ctx).
		Raw(query, scripId, userId, pq.Array(watchlistIds)).
		Scan(&result).Error

	return result, err
}

func (repo *watchlistsRepository) DeleteScripWithCTE(
	ctx context.Context,
	db *gorm.DB,
	userId uint64,
	watchlistIds []uint64,
	scripId string,
) ([]structModels.WatchlistWithId, error) {

	var result []structModels.WatchlistWithId

	query := `
    WITH deleted AS (
		DELETE FROM watchlist_scrips ws
		WHERE ws.watchlist_id = ANY(?)
			AND ws.scrip_id = ?
			AND ws.watchlist_id IN (
				SELECT id FROM watchlists WHERE user_id = ?
			)
		RETURNING ws.watchlist_id
	),
	updated AS (
		UPDATE watchlists w 
		SET scrip_count = scrip_count - 1,
			last_updated_at = NOW()
		FROM deleted d  
		WHERE w.id = d.watchlist_id
		RETURNING w.id, w.watchlist_name
	)
	SELECT id as watchlist_id, watchlist_name FROM updated;
    `

	err := db.WithContext(ctx).
		Raw(query, pq.Array(watchlistIds), scripId, userId).
		Scan(&result).Error

	return result, err
}
