package repository

import (
	"authentication/commons/constants"
	"authentication/models"
	"context"
	"time"

	genericModels "stock_broker_application/src/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SigninUserRepository interface {
	SigninNewUser(ctx context.Context, db *gorm.DB, bffSigninUserRequest models.BFFSigninUserRequest) (*genericModels.User, error)
}

type signinUserRepository struct{}

func NewSigninUserRepository() *signinUserRepository { //constructor
	return &signinUserRepository{}
}

func (user *signinUserRepository) SigninNewUser(ctx context.Context, db *gorm.DB, bffSigninUserRequest models.BFFSigninUserRequest) (*genericModels.User, error) {
	start := time.Now()
	logger := logrus.New()

	var ExistingUser genericModels.User

	// SELECT * FROM users WHERE username = 'Arijit' LIMIT 1;
	result := db.WithContext(ctx).Where(constants.UsernameField, bffSigninUserRequest.Username).First(&ExistingUser)

	if result.Error != nil {
		return nil, result.Error
	}

	logger.WithFields(logrus.Fields{
		"user":    bffSigninUserRequest.Username,
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserLoggedInSuccessMsg)

	return &ExistingUser, nil

}


