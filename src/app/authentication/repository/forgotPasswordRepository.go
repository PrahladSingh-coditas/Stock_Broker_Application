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

type ForgotPasswordRepository interface {
	ForgotPassword(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error)
	GenerateOTP(ctx context.Context, db *gorm.DB, username string, otp map[string]interface{}) error
}

type forgotPasswordRepository struct{}

func NewForgotPasswordRepository() *forgotPasswordRepository {
	return &forgotPasswordRepository{}
}

func (user *forgotPasswordRepository) ForgotPassword(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error) {

	start := time.Now()
	logger := logrus.New()

	var existingUser genericModels.User

	result := db.WithContext(ctx).
		Table(constants.UsersTableName).
		Where(constants.FieldUsername, username).
		First(&existingUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.UserNotFoundError)
		}
		return nil, fmt.Errorf("%s: %w", constants.AuthenticationFailedError, result.Error)
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserReadSuccessMsg)

	return &existingUser, nil
}

func (user *forgotPasswordRepository) GenerateOTP(ctx context.Context, db *gorm.DB, username string, otp map[string]interface{}) error {
	start := time.Now()
	logger := logrus.New()

	result := db.Model(&genericModels.User{}).
		Where(constants.FieldUsername, username).
		Updates(otp)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New(constants.NoRecordsAffectedError)
		}
		return fmt.Errorf("%s: %w", constants.AuthenticationFailedError, result.Error)
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.ForgotPasswordGenerateOtpSuccessMsg)

	return nil
}
