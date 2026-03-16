package repository

import (
	"authentication/commons/constants"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ChangePasswordRepository interface {
	UpdateUserPassword(ctx context.Context, db *gorm.DB, username string, NewHashedPassword string, logger *logrus.Logger) error
}

type changePasswordRepository struct{}

func NewChangePasswordRepository() *changePasswordRepository {
	return &changePasswordRepository{}
}

func (repo *changePasswordRepository) UpdateUserPassword(ctx context.Context, db *gorm.DB, username string, newHashedPassword string, logger *logrus.Logger) error {

	start := time.Now()

	result := db.WithContext(ctx).Table(constants.UsersTableName).Where(constants.UsernameCondition, username).Update(constants.Password, newHashedPassword)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	if result.Error != nil {
		return result.Error
	}

	logger.WithFields(logrus.Fields{
		constants.User:    username,
		constants.Latency: time.Since(start),
	}).Info(constants.PasswordUpdateMsg)

	return nil
}
