package repository

import (
	"authentication/commons/constants"
	"context"
	"errors"
	"fmt"
	GenericUserModel "stock_broker_application/src/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SignInUserRepository interface {
	SignInUser(ctx context.Context, username string) (*GenericUserModel.User, error)
	StoreOTP(ctx context.Context, username string, data map[string]interface{}) error
}

type signInUserRepository struct {
	gdb *gorm.DB
}

func NewSignInUserRepository(gormDB *gorm.DB) *signInUserRepository {
	return &signInUserRepository{
		gdb: gormDB,
	}
}

func (repo *signInUserRepository) SignInUser(ctx context.Context, username string) (*GenericUserModel.User, error) {

	start := time.Now()
	logger := logrus.New()

	var fetchedUserData GenericUserModel.User
	var user GenericUserModel.User

	err := repo.gdb.Model(&user).Where(constants.UsernameCondition, username).First(&fetchedUserData)

	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return nil, errors.New(constants.ErrUserNotFoundMsg)
		}
		return nil, errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	logger.WithFields(logrus.Fields{
		constants.User:    username,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info(constants.UserDataFetchedMsg)

	return &fetchedUserData, nil
}

func (repo *signInUserRepository) StoreOTP(ctx context.Context, username string, updates map[string]interface{}) error {

	var user GenericUserModel.User

	err := repo.gdb.Model(&user).
		Where(constants.UsernameCondition, username).
		Updates(updates)

	fmt.Println(err.Error)

	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return errors.New(constants.ErrUserNotFoundMsg)
		}
		return errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	return nil
}
