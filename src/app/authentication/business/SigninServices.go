package business

import (
	//constants "authentication/commons/constants"
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"

	//genericErrors "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
)

type SigninUserService struct {
	signinUserRepository repository.SigninUserRepository
}

func NewSigninUserService(signinUserRepository repository.SigninUserRepository) *SigninUserService {
	return &SigninUserService{
		signinUserRepository: signinUserRepository,
	}
}

func (service *SigninUserService) SigninUser(ctx context.Context, spanCtx context.Context, bffSigninUserRequest models.BFFSigninUserRequest) error {
	//returns struct of grom
	postgresClient := utils.GetPostgresClient()
	client := postgresClient.GormDB

	user, err := service.signinUserRepository.SigninNewUser(spanCtx, client, bffSigninUserRequest)
	if err != nil {
		return errors.New(constants.ErrInvalidEmailorPassword)
	}

	passwordMatch := utils.CompareHashPassword(user.Password, bffSigninUserRequest.Password)
	if !passwordMatch {
		return errors.New(constants.ErrPasswordMismatch)
	}

	return nil
}
