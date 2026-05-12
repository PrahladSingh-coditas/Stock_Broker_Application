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
	gormDB *gorm.DB
}

func NewValidateUserOtpRepository(gdb *gorm.DB) *validateUserOtpRepository {
	return &validateUserOtpRepository{gormDB: gdb}
}

func (repo *validateUserOtpRepository) GetUserByUsername(ctx context.Context, username string) (*genericModels.User, error) {

	start := time.Now()
	logger := logrus.New()
	var user genericModels.User

	result := repo.gormDB.WithContext(ctx).Table(constants.UsersTableName).Where(constants.UsernameCondition, username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	logger.WithFields(logrus.Fields{
		constants.User:    username,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info(constants.UserOtpFetchedMsg)

	return &user, nil
}
