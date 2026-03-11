package repository

import (
	"authentication/commons/constants"
	"context"
	GenericUserModel "stock_broker_application/src/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ChangePasswordReposioty interface {
	UpdateUserPassword(ctx context.Context, db *gorm.DB, username string, NewHashedPassword string) error
}

type changePasswordRepository struct{}

func NewChangePasswordRepository() *changePasswordRepository {
	return &changePasswordRepository{}
}

func (repo *changePasswordRepository) UpdateUserPassword(ctx context.Context, db *gorm.DB, username string, NewHashedPassword string) error {

	start := time.Now()
	logger := logrus.New()

	var user GenericUserModel.User
	result := db.Model(&user).Where(constants.UsernameCondtion, username).Update(constants.Password, NewHashedPassword)

	if result.Error != nil {
		return result.Error
	}

	logger.WithFields(logrus.Fields{
		constants.User:    username,
		constants.Latency: time.Since(start),
	}).Info(constants.PasswordUpdateMsg)

	return nil
}
