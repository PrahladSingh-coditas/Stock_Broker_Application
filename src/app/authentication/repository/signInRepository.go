package repository

import (
	"authentication/commons/constants"
	"context"
	"errors"
	genericModels "stock_broker_application/src/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SigninUserRepository interface {
	SigninUser(ctx context.Context, username string) (*genericModels.User, error)
}

type signinUserRepository struct {
	db *gorm.DB
}

func NewSigninUserRepository(gdb *gorm.DB) *signinUserRepository {
	return &signinUserRepository{db: gdb}
}

func (user *signinUserRepository) SigninUser(ctx context.Context, username string) (*genericModels.User, error) {

	start := time.Now()
	logger := logrus.New()

	var existingUser genericModels.User

	result := user.db.WithContext(ctx).
		Table(constants.UsersTableName).
		Where(constants.FieldUsername, username).
		First(&existingUser)

	//we have checked here with username
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.UserNotFoundError)
		}
		return nil, result.Error
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserLoggedInSuccessMsg)

	return &existingUser, nil
}
