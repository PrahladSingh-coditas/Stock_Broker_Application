package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"
)

type SigninUserService struct {
	signinUserRepository repository.SigninUserRepository // we have taken a field that is an interface in repository
}

func NewSigninUserService(signinUserRepository repository.SigninUserRepository) *SigninUserService {
	return &SigninUserService{
		signinUserRepository: signinUserRepository, //constructor
	}
}

func (service *SigninUserService) SigninUser(ctx context.Context, spanCtx context.Context, bffSigninUserRequest models.BFFSigninUserRequest) error {
	// postgresClinet := utils.GetPostgresClient()
	// client := postgresClinet.GormDB
	user, err := service.signinUserRepository.SigninUser(spanCtx, bffSigninUserRequest.Username)
	if err != nil {
		if err.Error() == constants.UserNotFoundError {
			return errors.New(constants.UserNotFoundError)
		}
		return err
	}

	passwordMatch := utils.CompareHashPassword(user.Password, bffSigninUserRequest.Password)
	if !passwordMatch {
		return errors.New(constants.InvalidUsernamePasswordError)
	}
	return nil

}
