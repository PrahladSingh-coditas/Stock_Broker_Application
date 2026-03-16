package business

import (
	"authentication/commons/constants"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ChangePasswordService struct {
	changePasswordReposioty repository.ChangePasswordRepository
}

func NewChangePasswordService(changePasswordReposioty repository.ChangePasswordRepository) *ChangePasswordService {
	return &ChangePasswordService{
		changePasswordReposioty: changePasswordReposioty,
	}
}

func (service *ChangePasswordService) ServiceChangePassword(ctx context.Context, spanCtx context.Context, username string, newPassword string, logger *logrus.Logger) error {

	postgresClinet := utils.GetPostgresClient()
	tx := postgresClinet.GormDB

	newHashedPassword, err := utils.HashPassword(newPassword)

	if err != nil {
		return errors.New(constants.ErrFailedToEncrypt)
	}

	err = service.changePasswordReposioty.UpdateUserPassword(ctx, tx, username, newHashedPassword, logger)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(constants.ErrUserNotFoundMsg)
		}
		return errors.New(constants.ErrDatabaseQueryErrorMsg)
	}

	return nil
}
