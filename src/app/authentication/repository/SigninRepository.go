package repository

import (
	"authentication/commons/constants"
	"context"
	"errors"

	genericModels "stock_broker_application/src/models"

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
	var user genericModels.User

	result := db.WithContext(ctx).Table(constants.UsersTableName).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil,
				errors.New(constants.ErrUsernameNotFound)
		}
		return nil, result.Error
	}
	return &user, nil
}
