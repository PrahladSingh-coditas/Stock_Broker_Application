package repository

import (
	"context"
	"errors"
	genericModels "stock_broker_application/src/models"
	"strings"
	"watchlists/commons/constants"
	"watchlists/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AdgStoWatchlistRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*genericModels.User, error)
	GetWatchlistDetails(ctx context.Context, addedTo []uint64) ([]genericModels.Watchlist, error)
	DelScripFromWatchlists(ctx context.Context, userId uint64, scripId string, watchlistIds []uint64) ([]uint64, error)
	AddScripToWatchlists(ctx context.Context, userId uint64, scripId string, watchlistIds []uint64) ([]uint64, []uint64, error)
	GetScripFromWatchlists(ctx context.Context, userId uint64, scripId string) ([]models.WatchlistWithId, error)
}

type AdgStoWatchlistsRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewadgStoWatchlistsRepository(db *gorm.DB, redisClient *redis.Client) *AdgStoWatchlistsRepository {
	return &AdgStoWatchlistsRepository{
		db:          db,
		redisClient: redisClient,
	}
}

// checks if the user exists in the given table or not
func (repo *AdgStoWatchlistsRepository) GetUserByUsername(ctx context.Context, username string) (*genericModels.User, error) {
	var user genericModels.User

	result := repo.db.WithContext(ctx).
		Table(constants.UsersTableName).
		Where(constants.Username, username).
		First(&user) //change table name
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// returns a list of watchlistWithIds for returning in the service response
// optimised as compared to other helper functions, because here instaed of using a loop, we used IN oprator
func (repo *AdgStoWatchlistsRepository) GetWatchlistDetails(ctx context.Context, addedTo []uint64) ([]genericModels.Watchlist, error) {
	var watchlists []genericModels.Watchlist

	err := repo.db.WithContext(ctx).
		Table("watchlists").
		Where("id IN ?", addedTo).
		Find(&watchlists).Error

	if err != nil {
		return nil, constants.DatabaseQueryError
	}

	return watchlists, nil
}

// deletes the scripid in the request from the mentioned watchlists,
// returns a []uint64 for the list of WIds which belong to the user.
func (repo *AdgStoWatchlistsRepository) DelScripFromWatchlists(ctx context.Context, userId uint64, scripId string, watchlistIds []uint64) ([]uint64, error) {
	if len(watchlistIds) == 0 {
		return []uint64{}, nil
	}

	query := `
	WITH valid_watchlists AS(
		SELECT id 
		FROM watchlists 
		WHERE user_id = ? 
		AND id IN (?)
	),
	deleted_from AS (
		DELETE FROM watchlist_scrips
		WHERE scrip_id = ?
		AND watchlist_id IN (SELECT id FROM valid_watchlists)
		RETURNING watchlist_id
	),
	update_count AS (
		UPDATE watchlists
		SET scrip_count = GREATEST(scrip_count - 1, 0),
			last_updated_at = NOW()
		WHERE id IN (SELECT watchlist_id FROM deleted_from)
	)
	SELECT id FROM valid_watchlists
	`
	var validWatchlists []uint64
	err := repo.db.WithContext(ctx).Raw(query, userId, watchlistIds, scripId).Scan(&validWatchlists).Error
	if err != nil {
		return nil, constants.DatabaseQueryError
	}

	return validWatchlists, nil
}

func (repo *AdgStoWatchlistsRepository) AddScripToWatchlists(ctx context.Context, userId uint64, scripId string, watchlistIds []uint64) ([]uint64, []uint64, error) {

	type Result struct {
		ValidWatchlistIds    uint64  `gorm:"column:valid_watchlist_ids"`
		InsertedWatchlistIds *uint64 `gorm:"column:inserted_watchlist_ids"` //ptr because all vw.id will appear, but ir.watchlist_id maybe null, go has no null for uint64,
		// so when took ptr, corresponding value for null becomes nil in go
	}

	results := []Result{}
	valid := []uint64{}
	inserted := []uint64{}

	query := `
	WITH valid_watchlists AS (
			SELECT id, scrip_count
			FROM watchlists w
			WHERE w.id IN (?)
			AND w.user_id = ?
		),
		insert_records AS (
			INSERT INTO watchlist_scrips (watchlist_id, scrip_id)
			SELECT vw.id, ?
			FROM valid_watchlists vw
			WHERE vw.scrip_count < 10
			ON CONFLICT (watchlist_id, scrip_id) DO NOTHING
			RETURNING watchlist_id  
		),
		update_count AS (
			UPDATE watchlists
			SET scrip_count = scrip_count + 1,
				last_updated_at = NOW()
			WHERE id IN (SELECT watchlist_id FROM insert_records)
		)
		SELECT vw.id AS valid_watchlist_ids,
			ir.watchlist_id AS inserted_watchlist_ids
		FROM valid_watchlists AS vw
		LEFT JOIN insert_records AS ir
			ON vw.id = ir.watchlist_id;
	`
	err := repo.db.WithContext(ctx).Raw(query, watchlistIds, userId, scripId).Scan(&results).Error
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "foreign key") {
			return nil, nil, constants.ScripDoesNotExistError
		}
		return nil, nil, err
	}

	for _, res := range results {
		valid = append(valid, res.ValidWatchlistIds)
		if res.InsertedWatchlistIds != nil {
			inserted = append(inserted, *res.InsertedWatchlistIds)
		}
	}

	return valid, inserted, nil
}

// checks if the scrip exists in scrip_masters table and if exists, returns the list of watchlists in which the scrip is present,
// returns a []uint64 for the list of WIds in which the scrip is present
func (repo *AdgStoWatchlistsRepository) GetScripFromWatchlists(ctx context.Context, userId uint64, scripId string) ([]models.WatchlistWithId, error) {
	var watchlists []models.WatchlistWithId
	type result struct {
		ScripCount    int     `gorm:"column:scrip_count"`
		Id            *uint64 `gorm:"column:id"`
		WatchlistName *string `gorm:"column:watchlist_name"`
	}
	var resultRows []result

	query := `
		WITH valid_scrip AS (
			SELECT COUNT(*) AS scrip_count
			FROM scrip_masters
			WHERE id = ?
		),
		filtered AS (
			SELECT w.id, w.watchlist_name
			FROM watchlists w
			JOIN watchlist_scrips ws
				ON w.id = ws.watchlist_id
			WHERE w.user_id = ?
				AND ws.scrip_id = ?
		)
		SELECT
			vs.scrip_count,
			f.id,
			f.watchlist_name
		FROM valid_scrip vs
		LEFT JOIN filtered f
			ON TRUE
	`
	err := repo.db.WithContext(ctx).Raw(query, scripId, userId, scripId).Scan(&resultRows).Error
	if err != nil {
		return nil, constants.DatabaseQueryError
	}

	if len(resultRows) == 0 {
		return nil, constants.DatabaseQueryError
	}

	if resultRows[0].ScripCount == 0 {
		return nil, constants.ScripDoesNotExistError
	}

	for _, res := range resultRows {
		if res.Id != nil {
			watchlists = append(watchlists, models.WatchlistWithId{
				Id:   *res.Id,
				Name: *res.WatchlistName,
			})
		}
	}

	return watchlists, nil
}
