package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"watchlist/commons/constants"
	"watchlist/models"

	genericModels "stock_broker_application/src/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/lib/pq"
)

type WatchlistRepository interface {
	GetUserFromDb(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error)
	WatchlistGetOperation(ctx context.Context, db *gorm.DB, UserId uint64, ScripId string) ([]models.WatchlistWithId, error)
	WatchlistAddOperation(ctx context.Context, db *gorm.DB, UserId uint64, ScripId string, WatchlistIds []uint64) ([]models.WatchlistWithId, []models.WatchlistValidation, error)
	WatchlistDeleteOperation(ctx context.Context, db *gorm.DB, UserId uint64, ScripId string, WatchlistIds []uint64) ([]models.WatchlistWithId, []models.WatchlistValidation, error)
}

type watchlistRepository struct{}

func NewWatchlistRepository() *watchlistRepository {
	return &watchlistRepository{}
}

func (user *watchlistRepository) GetUserFromDb(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error) {
	start := time.Now()
	logger := logrus.New()

	var ExistingUser genericModels.User

	result := db.WithContext(ctx).Table("users").
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

func (user *watchlistRepository) WatchlistGetOperation(ctx context.Context, db *gorm.DB, UserId uint64, ScripId string) ([]models.WatchlistWithId, error) {
	start := time.Now()
	logger := logrus.New()

	var WatchlistNamewithId []models.WatchlistWithId

	result := db.WithContext(ctx).
		Table("watchlists w").
		Select(" DISTINCT w.id AS watchlist_id, w.watchlist_name as watchlist_name").
		Joins("JOIN watchlist_scrips ws ON w.id = ws.watchlist_id").
		Where("w.user_id= ? and ws.scrip_id=?", UserId, ScripId).
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
func (repo *watchlistRepository) ValidateWatchlistOperation(ctx context.Context, db *gorm.DB, UserId uint64, ScripId string, WatchlistIds []uint64) ([]models.WatchlistValidation, error) {

	var results []models.WatchlistValidation

	query := `
	WITH input_watchlists AS (
		SELECT UNNEST(?::bigint[]) AS watchlist_id
	),

	scrip_exists AS (
		SELECT COUNT(*) > 0 AS exists_flag
		FROM scrip_masters
		WHERE id = ?
	),

	watchlist_user_check AS (
		SELECT id AS watchlist_id, user_id
		FROM watchlists
		WHERE id IN (SELECT watchlist_id FROM input_watchlists)
	),

	scrip_count_check AS (
		SELECT watchlist_id,
			COUNT(*) AS total_scrips
		FROM watchlist_scrips
		GROUP BY watchlist_id
	),

	scrip_present_check AS (
		SELECT watchlist_id
		FROM watchlist_scrips
		WHERE LOWER(scrip_id) = LOWER(?)
		AND watchlist_id IN (SELECT watchlist_id FROM input_watchlists)
	)

	SELECT
		i.watchlist_id AS "WatchlistId",

		-- ADD
		CASE
			WHEN NOT (SELECT exists_flag FROM scrip_exists) THEN false
			WHEN wu.watchlist_id IS NULL THEN false
			WHEN wu.user_id != ? THEN false
			WHEN COALESCE(sc.total_scrips,0) >= 10 THEN false
			WHEN sp.watchlist_id IS NOT NULL THEN false
			ELSE true
		END AS "CanInsert",

		-- DELETE
		CASE
			WHEN wu.watchlist_id IS NULL THEN false
			WHEN wu.user_id != ? THEN false
			WHEN sp.watchlist_id IS NULL THEN false
			ELSE true
		END AS "CanDelete",

		CASE
			WHEN wu.watchlist_id IS NULL
				THEN 'WATCHLIST_NOT_FOUND'
			WHEN wu.user_id != ?
				THEN 'WATCHLIST_NOT_OF_USER'
			WHEN sp.watchlist_id IS NOT NULL
				AND COALESCE(sc.total_scrips,0) < 10
				THEN 'SCRIP_ALREADY_PRESENT(Duplicate)'
			WHEN sp.watchlist_id IS NULL
				THEN 'SCRIP_NOT_PRESENT'
			ELSE 'OK'
		END AS "Reason"

	FROM input_watchlists i
	LEFT JOIN watchlist_user_check wu
		ON wu.watchlist_id = i.watchlist_id
	LEFT JOIN scrip_count_check sc
		ON sc.watchlist_id = i.watchlist_id
	LEFT JOIN scrip_present_check sp
		ON sp.watchlist_id = i.watchlist_id;
	`

	err := db.Raw(
		query,
		pq.Array(WatchlistIds),  // UNNEST
		ScripId,      // scrip_exists
		ScripId,      // scrip_present_check
		UserId,       // CanInsert ownership
		UserId,       // CanDelete ownership
		UserId,       // Reason ownership
	).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *watchlistRepository) WatchlistAddOperation(ctx context.Context, db *gorm.DB, UserId uint64, ScripId string, WatchlistIds []uint64) ([]models.WatchlistWithId, []models.WatchlistValidation, error) {

	start := time.Now()
	logger := logrus.New()

	validationResults, err := repo.ValidateWatchlistOperation(ctx, db, UserId, ScripId, WatchlistIds)
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
	logrus.Infof("Validation Results: %+v", validationResults)
	// append valid watchlists and their scipr ids to slice of struct
	var records []models.WatchlistAndScrip
	for _, id := range validWatchlists {
		records = append(records, models.WatchlistAndScrip{
			WatchlistId: id,
			ScripId:     ScripId,
		})
	}

	query := `
		WITH inserted AS (
			INSERT INTO watchlist_scrips (watchlist_id, scrip_id)
			SELECT id, ?
			FROM watchlists
			WHERE id IN (?)
			RETURNING watchlist_id
		)

		UPDATE watchlists w
		SET
			scrip_count = w.scrip_count + 1,
			last_updated_at = NOW()
		FROM inserted i
		WHERE w.id = i.watchlist_id
		RETURNING w.id AS "WatchlistId", w.watchlist_name AS "WatchlistName";
	`

	var result []models.WatchlistWithId
	err = db.WithContext(ctx).Raw(query, ScripId, validWatchlists).Scan(&result).Error

	if err != nil {
		return nil, nil, err
	}

	logger.WithFields(logrus.Fields{
		"scripId": ScripId,
		"latency": time.Since(start).Milliseconds(),
	}).Info("Watchlist ADD completed")

	return result, warnings, nil
}

func (repo *watchlistRepository) WatchlistDeleteOperation(ctx context.Context, db *gorm.DB, UserId uint64, ScripId string, WatchlistIds []uint64) ([]models.WatchlistWithId, []models.WatchlistValidation, error) {
	start := time.Now()
	logger := logrus.New()

	validations, err := repo.ValidateWatchlistOperation(ctx, db, UserId, ScripId, WatchlistIds)
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

	fmt.Printf("validations: %+v\n", validations)
	var deletedWatchlists []models.WatchlistWithId
	query := `
	WITH delete_from_watchlist_scrips AS(
			delete from watchlist_scrips
			where watchlist_id in (?) and scrip_id = ?
			returning watchlist_id
	),
	update_watchlists_table as(
			update watchlists
			set
			scrip_count = Greatest(scrip_count-1,0),
			last_updated_at= now()
			where id in (select watchlist_id from delete_from_watchlist_scrips)
	)
	SELECT 
		w.id AS watchlist_id,
		w.watchlist_name
	FROM watchlists w
	JOIN delete_from_watchlist_scrips d ON w.id = d.watchlist_id;
	`

	err = db.WithContext(ctx).Raw(query, validWatchlistIds, ScripId).Scan(&deletedWatchlists).Error
	if err != nil {
		return nil, nil, err
	}

	logger.WithFields(logrus.Fields{
		"scripId": ScripId,
		"latency": time.Since(start).Milliseconds(),
	}).Info("Watchlist DEL completed")

	return deletedWatchlists, warnings, nil

}

// func (user *watchlistRepository) ValidateScripIdExists(ctx context.Context, db *gorm.DB, ScripId string) (bool, error) {
// 	start := time.Now()
// 	logger := logrus.New()

// 	var count int64

// 	result := db.WithContext(ctx).
// 		Table("scrip_masters sm").
// 		Where("scrip_id=?", ScripId).
// 		Scan(&count)

// 	if result.Error != nil {
// 		return false, result.Error
// 	}
// 	logger.WithFields(logrus.Fields{
// 		"user":    ScripId,
// 		"latency": time.Since(start).Milliseconds(),
// 	}).Info("Validation of scrip id exists in master table success")

// 	return count > 0, nil
// }

// func (user *watchlistRepository) ValidateWatchlistOfUser(ctx context.Context, db *gorm.DB, UserId uint64, WatchlistIds []uint64) (bool, error) {
// 	start := time.Now()
// 	logger := logrus.New()

// 	var count uint64

// 	result := db.WithContext(ctx).
// 		Table("watchlists").
// 		Where("user_id=? and id IN ?", UserId, WatchlistIds).
// 		Scan(&count)

// 	if result.Error != nil {
// 		return false, result.Error
// 	}
// 	logger.WithFields(logrus.Fields{
// 		"user":    UserId,
// 		"latency": time.Since(start).Milliseconds(),
// 	}).Info("Validation of watchlist ids belong to current user success")

// 	return count == uint64(len(WatchlistIds)), nil
// }

// func (user *watchlistRepository) CheckScripLimit(ctx context.Context, db *gorm.DB, WatchlistIds []uint64) (bool, error) {
// 	start := time.Now()
// 	logger := logrus.New()

// 	type WatchlistIdWithScripCount struct {
// 		watchlistId uint64
// 		count       int64
// 	}

// 	var watchlistIdWithScripCount []WatchlistIdWithScripCount

// 	result := db.WithContext(ctx).
// 		Table("watchlist_scrips").
// 		Select("watchlist_id,COUNT(scrip_id) as count").
// 		Where("watchlist_id IN ?", WatchlistIds).
// 		Group("watchlist_id").
// 		Scan(&watchlistIdWithScripCount)

// 	if result.Error != nil {
// 		return false, result.Error
// 	}

// 	countMap := make(map[uint64]int64)

// 	for _, val := range watchlistIdWithScripCount {
// 		countMap[val.watchlistId] = val.count
// 	}

// 	for _, id := range WatchlistIds {
// 		if countMap[id] >= 10 {
// 			return false, errors.New("scrip limit exceeded")
// 		}
// 	}

// 	logger.WithFields(logrus.Fields{
// 		"user":    watchlistIdWithScripCount,
// 		"latency": time.Since(start).Milliseconds(),
// 	}).Info("Validation of scrip limit success")
// 	return true, nil
// }

// func (user *watchlistRepository) CheckScripExistsInWatchlists(ctx context.Context, db *gorm.DB,ScripId string, WatchlistIds []uint64) (bool,error) {
// 	start := time.Now()
// 	logger := logrus.New()

// 	var count int64

// 	result := db.WithContext(ctx).
// 	Table("watchlist_scrips").
// 	Where("scrip_id=? AND watchlist_id IN ?",ScripId, WatchlistIds).
// 	Scan(&count)

// 	if result.Error !=nil{
// 		return false, result.Error
// 	}

// 	if count>0{
// 		return false, errors.New("Scrip already exists in watchlist")
// 	}

// 	logger.WithFields(logrus.Fields{
// 		"user":    ScripId,
// 		"latency": time.Since(start).Milliseconds(),
// 	}).Info("Validation of scrip exist in watchlist success")

// 	return true,nil
// }

// func (user *watchlistRepository) WatchlistAddOperation(ctx context.Context, db *gorm.DB, UserId uint64, ScripId string, WatchlistIds []uint64) (bool,error) {
// 	start := time.Now()
// 	logger := logrus.New()

// 	scripIdExistsInScripMaster,err := user.ValidateScripIdExists(ctx,db,ScripId)
// 	if err != nil{
// 		return false,err
// 	}

// 	watchlistOfUser, err := user.ValidateWatchlistOfUser(ctx,db,UserId, WatchlistIds)
// 	if err!=nil{
// 		return false,err
// 	}

// 	scripLimit, err := user.CheckScripLimit(ctx,db,WatchlistIds)
// 	if err!=nil {
// 		return false,err
// 	}

// 	scripExists, err := user.CheckScripExistsInWatchlists(ctx,db,ScripId,WatchlistIds)
// 	if err!=nil {
// 		return false,err
// 	}

// 	if scripIdExistsInScripMaster==true && watchlistOfUser==true && scripLimit==true && scripExists==true{
// 		result := db.WithContext(ctx).
// 		Table("watchlists").
// 		Insert("INSERT into watchlist_scrips").
// 		Values("values watchlist_id=? AND scrip_id=?",WatchlistIds,ScripId)

// 		if result.Error != nil{
// 			return false, errors.New("DB Query Error")
// 		}
// 	}

// 	logger.WithFields(logrus.Fields{
// 		"user":    ScripId,
// 		"latency": time.Since(start).Milliseconds(),
// 	}).Info("Scrip and watchlist Id added successfully")

// 	return true,nil

// }
