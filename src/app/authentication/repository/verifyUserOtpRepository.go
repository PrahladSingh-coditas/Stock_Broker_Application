package repository

import (
	"authentication/commons/constants"
	"context"
	"fmt"
	genericModels "stock_broker_application/src/models"

	"gorm.io/gorm"
)

type ValidateUserOtpRepository interface {
	GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error)
}

type validateUserOtpRepository struct{}

func NewValidateUserOtpRepository() *validateUserOtpRepository {
	return &validateUserOtpRepository{}
}

func (repo *validateUserOtpRepository) GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (*genericModels.User, error) {
	var user genericModels.User
	result := db.WithContext(ctx).Table(constants.UsersTableName).Where(constants.Username, username).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf(constants.ErrUserNotFound)
	}
	return &user, nil
}
