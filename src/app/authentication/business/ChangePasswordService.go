package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"

	"github.com/pingcap/log"
)

type ChangePasswordService struct {
	changePasswordRepository repository.ChangePasswordRepository
}

func NewChangePasswordService(changePasswordRepository repository.ChangePasswordRepository) *ChangePasswordService {
	return &ChangePasswordService{
		changePasswordRepository: changePasswordRepository,
	}
}

func (service *ChangePasswordService) ChangePassword(ctx context.Context, spanCtx context.Context, bffChangePasswordRequest models.BFFChangePasswordRequest, username string) error {
	postgresClinet := utils.GetPostgresClient()
	client := postgresClinet.GormDB

	hashPassword, err := utils.HashPassword(bffChangePasswordRequest.NewPassword)																													
	if err != nil {
		log.Info(constants.ErrFailedToEncrypt)
	}

	new_password := hashPassword

	errs := service.changePasswordRepository.CheckUser(spanCtx, client, username, new_password)
	if errs != nil {
		if errs.Error() == constants.NoRecordsAffectedError {
			return errors.New(constants.NoRecordsAffectedError)
		}
		return errors.New(constants.AuthenticationFailedError)
	}

	return nil

}
