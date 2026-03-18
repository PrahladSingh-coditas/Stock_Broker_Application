package repository

import (
	"authentication/commons/constants"
	"context"
	"errors"
	"time"

	genericModels "stock_broker_application/src/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SignInRepository interface {
	GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error)
}

type signInRepository struct{}

func NewSignInRepository() *signInRepository {
	return &signInRepository{}
}

func (repo *signInRepository) GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error) {
	start := time.Now()
	logger := logrus.New()

	var user genericModels.User

	result := db.WithContext(ctx).Table(constants.UsersTableName).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil,
				errors.New(constants.ErrUsernameNotFound)
		}
		return nil, result.Error
	}
	
	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Microseconds(),
	}).Info(constants.UserLoggedInSuccessMsg)

	return &user, nil
}
