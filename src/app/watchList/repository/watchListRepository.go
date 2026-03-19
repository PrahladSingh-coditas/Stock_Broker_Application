package repository

import (
	"errors"
	"stock_broker_application/src/models"

	"gorm.io/gorm"
)

type WatchListRepository interface {
	GetUserByUsername(username string) (*models.User, error)

	IsScripExists(scripID string) (bool, error)

	AddScripToWatchlist(watchlistID uint64, scripID string) error

	GetWatchlistsByIDs(userID uint64, watchlistIDs []uint64) ([]models.Watchlist, error)

	IsScripInWatchlist(watchlistID uint64, scripID string) (bool, error)
}

type watchListRepository struct {
	db *gorm.DB
}

func NewWatchlistRepository(db *gorm.DB) *watchListRepository {
	return &watchListRepository{
		db: db,
	}
}

func (repo *watchListRepository) GetUserByUsername(username string) (*models.User, error) {

	var user models.User

	err := repo.db.
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *watchListRepository) AddScripToWatchlist(watchlistID uint64, scripID string) error {

	tx := repo.db.Begin()

	res := tx.Table("watchlists").
		Where("id = ? AND scrip_count < 10", watchlistID).
		Update("scrip_count", gorm.Expr("scrip_count + ?", 1))

	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("watchlist limit reached")
	}

	err := tx.Table("watchlist_scrips").Create(&models.WatchlistScrips{
		WatchlistID: watchlistID,
		ScripID:     scripID,
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (repo *watchListRepository) GetWatchlistsByIDs(userID uint64, watchlistIDs []uint64) ([]models.Watchlist, error) {

	var watchlists []models.Watchlist

	err := repo.db.Debug().
		Where("user_id = ? AND id IN ?", userID, watchlistIDs).
		Find(&watchlists).Error

	//userID not belong to user or no watchlist found for given ids

	return watchlists, err
}

func (repo *watchListRepository) IsScripInWatchlist(watchlistID uint64, scripID string) (bool, error) {
	var count int64

	err := repo.db.
		Model(&models.WatchlistScrips{}).
		Where("watchlist_id = ? AND scrip_id = ?", watchlistID, scripID).Debug().
		Count(&count).Error

	return count > 0, err
}

func (repo *watchListRepository) IsScripExists(scripID string) (bool, error) {

	var count int64

	err := repo.db.
		Model(&models.ScripMaster{}).
		Where("id = ?", scripID).Debug().
		Count(&count).Error

	return count > 0, err
}
