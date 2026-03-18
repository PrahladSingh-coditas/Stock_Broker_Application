package business

import (
	//constants "authentication/commons/constants"
	"authentication/commons/constants"
	"authentication/repository"
	"context"
	"errors"

	//genericErrors "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"

	"gorm.io/gorm"
)

type ChangePasswordService struct {
	changePasswordRepository repository.ChangePasswordRepository
}

func NewChangePasswordService(changePasswordRepository repository.ChangePasswordRepository) *ChangePasswordService {
	return &ChangePasswordService{
		changePasswordRepository: changePasswordRepository,
	}
}

func (service *ChangePasswordService) ChangePassword(ctx context.Context, spanCtx context.Context, username string, newPassword string, confirmPassword string) error {
	//returns struct of gorm
	postgresClient := utils.GetPostgresClient()
	client := postgresClient.GormDB

	if newPassword != confirmPassword{
		return errors.New(constants.ErrPasswordMismatch)
	}
	
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New(constants.ErrFailedToEncrypt)
	}

	err = service.changePasswordRepository.ChangePassword(spanCtx, client, username, hashedPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNoRowsAffected)
		}
		return err
	}

	return nil
}
