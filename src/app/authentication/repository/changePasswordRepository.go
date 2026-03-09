package repository

import (
	"stock_broker_application/src/models"
	"stock_broker_application/src/utils"
)

type ChangePasswordRepository interface {
	UpdatePassword(username string, Password string) error
}

type changePasswordRepository struct{}

func NewChangePasswordRepository() *changePasswordRepository {
	return &changePasswordRepository{}
}

func (repo *changePasswordRepository) UpdatePassword(username string, Password string) error {
	result := utils.GetPostgresClient().
		GormDB.Model(&models.User{}).Where("username = ?", username).
		Update("password", Password)
	return result.Error
}
