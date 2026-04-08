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
	signinUserRepository repository.SigninUserRepository//a service ka object(with type as repository.SigninUserRepository) that holds the repository so service can call repo ke functions
}

func NewSigninUserService(signinUserRepository repository.SigninUserRepository) *SigninUserService { //funciton where the dependecny injection takes place
	return &SigninUserService{
		signinUserRepository: signinUserRepository,
	}
}

func (service *SigninUserService) SigninUser(ctx context.Context, spanCtx context.Context, bffSigninUserRequest models.BFFSigninUserRequest) error { //contains actual business logic
	//returns struct of gorm
	postgresClient := utils.GetPostgresClient() //contains DB connection
	client := postgresClient.GormDB //taking the actual GORM DB ka instance form client

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
