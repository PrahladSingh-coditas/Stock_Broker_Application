package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"fmt"
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
	userDB, errs := service.changePasswordRepository.GetPassword(spanCtx, username)
	if errs != nil {
		if errs.Error() == constants.UserNotFoundError {
			return errors.New(constants.UserNotFoundError)
		}
		return errors.New(constants.AuthenticationFailedError)
	}

	passwordMatch := utils.CompareHashPassword(userDB.Password, bffChangePasswordRequest.NewPassword)
	if !passwordMatch {
		hashPassword, err := utils.HashPassword(bffChangePasswordRequest.NewPassword)
		if err != nil {
			log.Info(constants.ErrFailedToEncrypt)
		}
		new_password := hashPassword
		err = service.changePasswordRepository.UpdatePassword(spanCtx, username, new_password)
		if err != nil {
			return fmt.Errorf(constants.PasswordChangeFailedError, err)
		}
	} else {
		return errors.New(constants.SamePasswordError)
	}

	return nil

}
