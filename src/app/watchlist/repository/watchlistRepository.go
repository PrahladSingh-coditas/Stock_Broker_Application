package repository

import (
	"context"
	"watchlist/models"

	genericModels "stock_broker_application/src/models"

	"gorm.io/gorm"
)

type WatchlistRepository interface {
	GetUserIDFromDb(ctx context.Context, db *gorm.DB, username string) (uint64, error)
	WatchlistGet(ctx context.Context, db *gorm.DB, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) (*genericModels.User, error)
}

type watchlistRepository struct{}

func NewWatchlistRepository() *watchlistRepository {
	return &watchlistRepository{}
}

// func (user *watchlistRepository) GetUserIDFromDb(ctx context.Context, db *gorm.DB, username string) (uint64, error) {
// 	start := time.Now()
// 	logger := logrus.New()

// 	result := db.WithContext(ctx).Table("users").
// 	Where(constants.UsernameField, username)

// }

// func (user *watchlistRepository) WatchlistGet(ctx context.Context, db *gorm.DB, bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest) (*genericModels.User, error) {
// 	start := time.Now()
// 	logger := logrus.New()

// 	var ExistingUser genericModels.User

// }
