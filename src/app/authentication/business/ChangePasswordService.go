package business

import (
	"authentication/commons/constants"
	"authentication/repository"
	"context"
	"stock_broker_application/src/utils"

	"gorm.io/gorm"
)

type ChangePasswordService struct {
	changePasswordReposioty repository.ChangePasswordReposioty
}

func NewChangePasswordService(changePasswordReposioty repository.ChangePasswordReposioty) *ChangePasswordService {
	return &ChangePasswordService{
		changePasswordReposioty: changePasswordReposioty,
	}
}

func (service *ChangePasswordService) ServiceChangePassword(ctx context.Context, spanCtx context.Context, username string, NewPassword string) error {

	postgresClinet := utils.GetPostgresClient()
	tx := postgresClinet.GormDB

	NewHashedPassword, err := utils.HashPassword(NewPassword)

	if err != nil {
		return constants.PasswordEncryptFailedError
	}

	errDuringUpdate := service.changePasswordReposioty.UpdateUserPassword(ctx, tx, username, NewHashedPassword)

	if errDuringUpdate != nil {
		if errDuringUpdate == gorm.ErrRecordNotFound {
			return constants.UserNotFoundError
		}
		return constants.DatabaseQueryError
	}

	return nil
}
