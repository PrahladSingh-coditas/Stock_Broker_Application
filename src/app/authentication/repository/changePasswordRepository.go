package repository

import (
	"authentication/commons/constants"
	"context"
	"errors"
	"fmt"
	genericModels "stock_broker_application/src/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ChangePasswordRepository interface {
	GetPassword(ctx context.Context, username string) (*genericModels.User, error)
	UpdatePassword(ctx context.Context, username string, password string) error
}

type changePasswordRepository struct {
	db *gorm.DB
}

func NewChangePasswordRepository(gdb *gorm.DB) *changePasswordRepository {
	return &changePasswordRepository{
		db: gdb,
	}
}

func (user *changePasswordRepository) GetPassword(ctx context.Context, username string) (*genericModels.User, error) {

	start := time.Now()
	logger := logrus.New()

	var User genericModels.User
	//get old password from db by username
	result := user.db.WithContext(ctx).
		Table(constants.UsersTableName).
		Where(constants.FieldUsername, username).
		First(&User)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.UserNotFoundError)
		}
		return nil, fmt.Errorf("%s: %w", constants.AuthenticationFailedError, result.Error)
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserReadSuccessMsg)

	return &User, nil
}

func (user *changePasswordRepository) UpdatePassword(ctx context.Context, username string, password string) error {

	start := time.Now()
	logger := logrus.New()

	//update password in database by new password
	result := user.db.WithContext(ctx).Table(constants.UsersTableName).
		Where(constants.FieldUsername, username).
		Update(constants.FieldPassword, password)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New(constants.UserNotFoundError)
		}
		return fmt.Errorf("%s: %w", constants.AuthenticationFailedError, result.Error)
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.PasswordChangeSuccessMsg)

	return nil
}
