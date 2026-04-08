package repository

import (
	"authentication/commons"
	"authentication/commons/constants"
	"context"
	"time"

	

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ChangePasswordRepository interface {
	ChangePassword(ctx context.Context, db *gorm.DB, username string, newPassword string)  error
}

type changePasswordRepository struct{}

func NewChangePasswordRepository() *changePasswordRepository { //constructor
	return &changePasswordRepository{}
}

func (user *changePasswordRepository) ChangePassword(ctx context.Context, db *gorm.DB, username string, newPassword string)  error {
	start := time.Now()
	logger := logrus.New()

	
	result := db.WithContext(ctx).Table("users").
		Where(constants.UsernameField, username).
		Update(commons.Password, newPassword)

	
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return result.Error
	}

	logger.WithFields(logrus.Fields{
		"user":    username,
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserLoggedInSuccessMsg)

	return  nil

}
