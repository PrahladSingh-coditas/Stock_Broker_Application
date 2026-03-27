package repository

import (
	"context"
	"errors"
	"stock_broker_application/src/app/watchlist/commons/constants"
	"stock_broker_application/src/app/watchlist/models"
	genericModels "stock_broker_application/src/models"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WatchlistRepository interface {
	GetWatchlistsWithId(ctx context.Context, db *gorm.DB, logger *logrus.Logger, userID uint64, scripId string) (*models.ResultListsForGET, error)
	GetUserIdByUsername(ctx context.Context, db *gorm.DB, logger *logrus.Logger, username string) (*uint64, error)
	DeleteScripsFromWatchlists(ctx context.Context, db *gorm.DB, logger *logrus.Logger, userID uint64, watchlistIds []uint64, scripId string) (*models.ResultListsForDEL, error)
	AddScripsToWatchlists(ctx context.Context, db *gorm.DB, logger *logrus.Logger, userID uint64, watchlistIds []uint64, scripId string) (*models.ResultListsForADD, error)
}

type watchlistRepository struct{}

func NewWatchlistRepository() *watchlistRepository {
	return &watchlistRepository{}
}

func (repo *watchlistRepository) GetUserIdByUsername(ctx context.Context, db *gorm.DB, logger *logrus.Logger, username string) (*uint64, error) {
	start := time.Now()

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

func (repo *watchlistRepository) GetWatchlistsWithId(ctx context.Context, db *gorm.DB, logger *logrus.Logger, userID uint64, scripId string) (*models.ResultListsForGET, error) {
	var resultListsForGET models.ResultListsForGET
	start := time.Now()

	query:=`
			WITH scrip_check AS (
			SELECT COUNT(*) AS cnt
			FROM scrip_masters
			WHERE id = ?
		),
			get_ids AS(
			SELECT w.id , w.watchlist_name
			FROM watchlists w
			JOIN watchlist_scrips ws ON w.id=ws.watchlist_id 
			WHERE w.user_id = ? AND ws.scrip_id = ?
		)
			SELECT
			ARRAY(SELECT id FROM get_ids) AS ids,
			ARRAY(SELECT watchlist_name FROM get_ids) AS names,
			(SELECT cnt FROM scrip_check) AS scrip_check
		`

	result := db.WithContext(ctx).Raw(query,scripId,userID,scripId).Row().Scan(pq.Array(&resultListsForGET.WatchlistId),pq.Array(&resultListsForGET.WatchlistName),&resultListsForGET.ScripCheckCount)

		

	if result != nil {
		return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.UserId:  userID,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info("Watchlists with given scripId retrieved successfully")

	return &resultListsForGET, nil
}

func (repo *watchlistRepository) DeleteScripsFromWatchlists(ctx context.Context, db *gorm.DB, logger *logrus.Logger, userID uint64, watchlistIds []uint64, scripId string) (*models.ResultListsForDEL, error) {
	var resultListsForDEL models.ResultListsForDEL

	start := time.Now()

	query := `
		WITH valid_ids AS (
			SELECT id
			FROM watchlists
			WHERE id =ANY(?) AND user_id = ?
		),
		deleteScrips AS (
			DELETE FROM watchlist_scrips
			WHERE watchlist_id IN (SELECT id FROM valid_ids) AND scrip_id =?
			RETURNING watchlist_id
		),
		updateScripCount AS (
			UPDATE watchlists
			SET scrip_count=scrip_count-1, last_updated=NOW()
			WHERE id IN (SELECT watchlist_id FROM deleteScrips)
		)
		SELECT 
		ARRAY(SELECT id FROM valid_ids) AS valid_ids,
		ARRAY(SELECT watchlist_id FROM deleteScrips) AS deleted_ids
	`

	result := db.WithContext(ctx).Raw(query, pq.Array(watchlistIds), userID, scripId).Row().
	Scan(pq.Array(&resultListsForDEL.ValidWatchlistIds),pq.Array(&resultListsForDEL.DeletedWatchlistIds))

	if result != nil {
		return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.UserId:  userID,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info("Watchlist IDs Deleted successfully")

	return &resultListsForDEL, nil
}

func (repo *watchlistRepository) AddScripsToWatchlists(ctx context.Context, db *gorm.DB, logger *logrus.Logger, userID uint64, watchlistIds []uint64, scripId string) (*models.ResultListsForADD, error) {
	start := time.Now()
	var resultWatchlistIds models.ResultListsForADD

	query := `
		WITH valid_watchlists AS (
			SELECT id, scrip_count
			FROM watchlists
			WHERE id = ANY(?) AND user_id = ?
		),
		scrip_check AS (
			SELECT COUNT(*) AS cnt
			FROM scrip_masters
			WHERE id = ?
		),
		classify AS (
			SELECT 
			vw.id AS watchlist_id,
			vw.scrip_count,
			CASE
			WHEN EXISTS (
				SELECT 1 FROM watchlist_scrips ws WHERE ws.watchlist_id = vw.id AND ws.scrip_id = ?
			) THEN 'skipped'
			WHEN vw.scrip_count >= ? THEN 'exceeded'
			ELSE 'eligible'
			END AS status
			FROM valid_watchlists vw
		),
		inserting AS (
			INSERT INTO watchlist_scrips (watchlist_id, scrip_id)
			SELECT watchlist_id, ?
			FROM classify 
			WHERE status='eligible'
			RETURNING watchlist_id
		),
		updating AS (
			UPDATE watchlists w
			SET scrip_count=scrip_count+1, last_updated=NOW()
			FROM inserting i
			WHERE w.id = i.watchlist_id
		)
		SELECT
		ARRAY(SELECT watchlist_id FROM classify WHERE status='eligible') as added_ids,
		ARRAY(SELECT watchlist_id FROM classify WHERE status='exceeded') as limit_exceeded_ids,
		ARRAY(SELECT watchlist_id FROM classify WHERE status='skipped') as skipped_ids,
		(SELECT cnt FROM scrip_check) as scrip_count
	`

	row := db.WithContext(ctx).Debug().Raw(query, pq.Array(watchlistIds), userID, scripId, scripId, constants.MaxScripsPerWatchlist, scripId).Row()

	result := row.Scan(pq.Array(&resultWatchlistIds.AddedWatchlistIds),
		pq.Array(&resultWatchlistIds.LimitExceededWatchlistIds),
		pq.Array(&resultWatchlistIds.SkippedWatchlistIds),
		&resultWatchlistIds.ScripCount)

	if row == nil || result != nil {
		return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.UserId:  userID,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info("Add operation completed")

	return &resultWatchlistIds, nil
}
