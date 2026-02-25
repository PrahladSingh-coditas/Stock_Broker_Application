package repository

import (
	"authentication/commons/constants"
	"authentication/models"
	"context"
	GenericUserModel "stock_broker_application/src/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SignInUserRepository interface {
	SignInUser(ctx context.Context, db *gorm.DB, bffSignInUserRequest models.BFFSignInUserRequest) (*GenericUserModel.User, error)
	StoreOTP(ctx context.Context, db *gorm.DB, username string, data map[string]interface{}) error
}

type signInUserRepository struct{}

func NewSignInUserRepository() *signInUserRepository {
	return &signInUserRepository{}
}

func (user *signInUserRepository) SignInUser(ctx context.Context, db *gorm.DB, bffSignInUserRequest models.BFFSignInUserRequest) (*GenericUserModel.User, error) {

	start := time.Now()
	logger := logrus.New()

	var fetchedUserData GenericUserModel.User

	findUserError := db.Where(constants.UsernameCondtion, bffSignInUserRequest.Username).First(&fetchedUserData).Error

	if findUserError != nil {
		return nil, findUserError
	}

	logger.WithFields(logrus.Fields{
		constants.User:    bffSignInUserRequest.Username,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info(constants.UserDataFetchedMsg)

	return &fetchedUserData, nil
}

func (repo *signInUserRepository) StoreOTP(ctx context.Context, db *gorm.DB, username string, data map[string]interface{}) error {

	var user GenericUserModel.User

	result := db.Model(&user).
		Where(constants.UsernameCondtion, username).
		Updates(data)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}
