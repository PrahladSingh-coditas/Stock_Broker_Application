package repository

import (
	"authentication/commons/constants"
	"context"
	"errors"
	"stock_broker_application/src/models"

	"gorm.io/gorm"
)

type ChangePasswordRepository interface {
	GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (*models.User, error)
	UpdatePassword(ctx context.Context, db *gorm.DB, username string, password string) error
}

type changePasswordRepository struct{}

func NewChangePasswordRepository() *changePasswordRepository {
	return &changePasswordRepository{}
}

// takes username in input and returns user struct if user is found, else returns error
func (repo *changePasswordRepository) GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (*models.User, error) {
	var userFromDB models.User
	result := db.WithContext(ctx).Table(constants.UsersTableName).Where(constants.Username, username).First(&userFromDB)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &userFromDB, nil
}

// takes username, password in input and updates the password of the user in  the table
func (repo *changePasswordRepository) UpdatePassword(ctx context.Context, db *gorm.DB, username string, password string) error {

	result := db.WithContext(ctx).Table(constants.UsersTableName).Where(constants.Username, username).UpdateColumn("password", password)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
