package repository

import (
	"authentication/commons/constants"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ChangePasswordRepository interface {
	UpdateUserPassword(ctx context.Context, username string, NewHashedPassword string, logger *logrus.Logger) error
}

type changePasswordRepository struct {
	gdb *gorm.DB
}

func NewChangePasswordRepository(gormDB *gorm.DB) *changePasswordRepository {
	return &changePasswordRepository{gdb: gormDB}
}

func (repo *changePasswordRepository) UpdateUserPassword(ctx context.Context, username string, newHashedPassword string, logger *logrus.Logger) error {

	start := time.Now()

	result := repo.gdb.Debug().WithContext(ctx).Table(constants.UsersTableName).Where(constants.UsernameCondition, username).Update(constants.Password, newHashedPassword)

	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	logger.WithFields(logrus.Fields{
		constants.User:    username,
		constants.Latency: time.Since(start),
	}).Info(constants.PasswordUpdateMsg)

	return nil
}
