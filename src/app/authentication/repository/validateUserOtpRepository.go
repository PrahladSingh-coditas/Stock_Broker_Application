package repository

import (
	"authentication/commons/constants"
	"context"
	genericModels "stock_broker_application/src/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ValidateUserOtpRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*genericModels.User, error)
}

type validateUserOtpRepository struct {
	db *gorm.DB
}

func NewValidateUserOtpRepository(gdb *gorm.DB) *validateUserOtpRepository {
	return &validateUserOtpRepository{db: gdb}
}

// this function takes username from request(from service), fetces the user from db and returns the user
func (repo *validateUserOtpRepository) GetUserByUsername(ctx context.Context, username string) (*genericModels.User, error) {

	start := time.Now()
	logger := logrus.New()

	var user genericModels.User
	result := repo.db.WithContext(ctx).Table(constants.UsersTableName).Where(constants.FieldUsername, username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	logger.WithFields(logrus.Fields{
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.OtpValidatedSuccessMsg)

	return &user, nil
}
