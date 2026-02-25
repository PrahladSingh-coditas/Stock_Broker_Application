package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"stock_broker_application/src/utils"
	"time"
)

type SignInUserService struct {
	signInUserRepository repository.SignInUserRepository
}

func NewSignInUserService(signInUserRepository repository.SignInUserRepository) *SignInUserService {
	return &SignInUserService{
		signInUserRepository: signInUserRepository,
	}
}

func (service *SignInUserService) SignInUser(ctx context.Context, spanCtx context.Context, bffSignInRequest models.BFFSignInUserRequest) error {
	postgresClinet := utils.GetPostgresClient()
	tx := postgresClinet.GormDB

	userDataFromDB, errorFromRepository := service.signInUserRepository.SignInUser(spanCtx, tx, bffSignInRequest)
	if errorFromRepository != nil {
		return errorFromRepository
	}

	checkPassword := utils.CompareHashPassword(userDataFromDB.Password, bffSignInRequest.Password)
	if !checkPassword {
		return fmt.Errorf(constants.ErrPasswordMismatch, errors.New(constants.ErrPasswordNotMatch))
	}

	generatedOTP := uint64(rand.Intn(9000) + 1000)
	data := map[string]interface{}{
		constants.OTPSent:       generatedOTP,
		constants.OTPExpiryTime: time.Now().Add(constants.OtpTimeLimit * time.Minute),
	}

	errInStoringOTP := service.signInUserRepository.StoreOTP(ctx, tx, bffSignInRequest.Username, data)
	if errInStoringOTP != nil {
		return errInStoringOTP
	}

	return nil
}
