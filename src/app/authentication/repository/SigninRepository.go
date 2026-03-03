package repository

import (
	"authentication/commons/constants"
	"context"
	"errors"
	GenericUserModel "stock_broker_application/src/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SignInUserRepository interface {
	SignInUser(ctx context.Context, db *gorm.DB, username string) (*GenericUserModel.User, error)
	StoreOTP(ctx context.Context, db *gorm.DB, username string, data map[string]interface{}) error
}

type signInUserRepository struct{}

func NewSignInUserRepository() *signInUserRepository {
	return &signInUserRepository{}
}

func (user *signInUserRepository) SignInUser(ctx context.Context, db *gorm.DB, username string) (*GenericUserModel.User, error) {

	start := time.Now()
	logger := logrus.New()

	var fetchedUserData GenericUserModel.User

	err := db.Where(constants.UsernameCondtion, username).First(&fetchedUserData)

	if err.RowsAffected == 0 {
		return nil, errors.New(constants.ErrUserNotFoundMsg)
	}

	if err.Error != nil {
		return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.User:    username,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info(constants.UserDataFetchedMsg)

	return &fetchedUserData, nil
}

func (repo *signInUserRepository) StoreOTP(ctx context.Context, db *gorm.DB, username string, updates map[string]interface{}) error {

	var user GenericUserModel.User

	result := db.Model(&user).
		Where(constants.UsernameCondtion, username).
		Updates(updates)

	if result.RowsAffected == 0 {
		return errors.New(constants.ErrUserNotFoundMsg)
	}

	if result.Error != nil {
		return errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	return nil
}
