package business

import (
	//constants "authentication/commons/constants"
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"math/rand/v2"

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
	//returns struct of gorm
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
 
	otp := 1000 + rand.Uint64N(9000)//prev range of Uint64N(9000) :[0,9000) now 1000 is added to entire range and the range becomes: [1000,10000)

	otpError := service.signinUserRepository.InsertOtpInDb(spanCtx, client, bffSigninUserRequest, otp)
	if otpError != nil {
		return errors.New(constants.ErrOtpFailed)
	}

	return nil
}
