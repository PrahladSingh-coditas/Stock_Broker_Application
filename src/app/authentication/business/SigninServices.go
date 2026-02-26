package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"

	"stock_broker_application/src/utils"
)

type SignInService struct {
	signInRepository repository.SignInRepository // Interface to Repository layer
}

func NewSignInService(SignInRepository repository.SignInRepository) *SignInService {
	return &SignInService{
		signInRepository: SignInRepository,
	}

}

func (service *SignInService) SignIn(ctx context.Context, bffSignInRequest models.BFFSignInRequest) (*models.BFFSignInResponse, error) {
	db := utils.GetPostgresClient().GormDB

	user, err := service.signInRepository.GetUserByUsername(ctx, db, bffSignInRequest.Username)
	if err != nil {
		return nil, err
	}

	// compare password
	if !utils.CompareHashPassword(user.Password, bffSignInRequest.Password) {
		return nil, errors.New(constants.ErrPasswordMismatch)
	}

	return &models.BFFSignInResponse{
		Message: constants.UserLoggedInSuccessMsg,
	}, nil

}
