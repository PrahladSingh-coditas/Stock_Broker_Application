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
	CheckUser(ctx context.Context, db *gorm.DB, username string, new_password string) error
}

type changePasswordRepository struct{}

func NewChangePasswordRepository() *changePasswordRepository {
	return &changePasswordRepository{}
}

func (user *changePasswordRepository) CheckUser(ctx context.Context, db *gorm.DB, username string, password string) error {

	start := time.Now()
	logger := logrus.New()

	//update password in database by new password
	result := db.Model(&genericModels.User{}).
		Where(constants.FieldUsername, username).
		Update(constants.FieldPassword, password)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New(constants.NoRecordsAffectedError)

		}
		return fmt.Errorf("%s: %w", constants.AuthenticationFailedError, result.Error)
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.PasswordChangeSuccessMsg)

	return nil
}
