package repository

import (
	"authentication/commons/constants"
	"authentication/models"
	"context"
	"errors"
	"time"

	genericModels "stock_broker_application/src/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SigninUserRepository interface {
	SigninAndInsertOtpInDb(ctx context.Context, db *gorm.DB, bffSigninUserRequest models.BFFSigninUserRequest, otp uint64) (*genericModels.User, error) 

}

type signinUserRepository struct{}

func NewSigninUserRepository() *signinUserRepository { //constructor
	return &signinUserRepository{}
}

func (user *signinUserRepository) SigninAndInsertOtpInDb(ctx context.Context, db *gorm.DB, bffSigninUserRequest models.BFFSigninUserRequest, otp uint64) (*genericModels.User, error) {
	start := time.Now()
	logger := logrus.New()

	var ExistingUser genericModels.User

	result := db.WithContext(ctx).
		Model(&genericModels.User{}).
		Where(constants.UsernameField, bffSigninUserRequest.Username).
		Clauses(clause.Returning{}).
		Updates(map[string]interface{}{
			"otpSent":      otp,
			"otpExpiresAt": time.Now().Unix() + 120,
		}).
		Scan(&ExistingUser)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("invalid username or password")
	}
	logger.WithFields(logrus.Fields{
		"user":    bffSigninUserRequest.Username,
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserLoggedInSuccessMsg)

	return &ExistingUser, nil
}
