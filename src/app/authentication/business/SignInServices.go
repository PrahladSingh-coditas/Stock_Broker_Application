package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
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

	userData, err1 := service.signInUserRepository.SignInUser(spanCtx, bffSignInRequest.Username)
	if err1 != nil {
		return err1
	}

	checkPassword := utils.CompareHashPassword(userData.Password, bffSignInRequest.Password)
	if !checkPassword {
		return errors.New(constants.ErrPasswordNotMatch)
	}

	generatedOTP := uint64(rand.Intn(9000) + 1000)
	data := map[string]interface{}{
		constants.OTPSent:       generatedOTP,
		constants.OTPExpiryTime: time.Now().Unix() + constants.OtpTimeLimitInSeconds,
	}

	err2 := service.signInUserRepository.StoreOTP(ctx, bffSignInRequest.Username, data)
	if err2 != nil {
		return err2
	}

	return nil
}
