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

type ForgotPasswordService struct {
	forgotPasswordRepository repository.ForgotPasswordRepository // we have taken a field that is an interface in repository
}

func NewForgotPasswordService(forgotPasswordRepository repository.ForgotPasswordRepository) *ForgotPasswordService {
	return &ForgotPasswordService{
		forgotPasswordRepository: forgotPasswordRepository, //constructor
	}
}

func (service *ForgotPasswordService) ReadRecordsWithConditions(ctx context.Context, spanCtx context.Context, bffForgotPasswordRequest models.BFFForgotPasswordRequest) error {
	postgresClinet := utils.GetPostgresClient()
	client := postgresClinet.GormDB

	user, err := service.forgotPasswordRepository.ForgotPassword(spanCtx, client, bffForgotPasswordRequest.Username)
	if err != nil {
		if err.Error() == constants.UserNotFoundError {
			return errors.New(constants.UserNotFoundError)
		}
		return errors.New(constants.AuthenticationFailedError)
	}

	if user.PanCard != bffForgotPasswordRequest.PanCard || user.PhoneNumber != bffForgotPasswordRequest.PhoneNumber {
		return errors.New(constants.AuthenticationFailedError)
	}

	otp := map[string]interface{}{
		"OtpSent":      rand.Intn(9000) + 1000,
		"OtpExpiresAt": uint64(time.Now().Unix() + 120),
	}

	errs := service.forgotPasswordRepository.GenerateOTP(spanCtx, client, bffForgotPasswordRequest.Username, otp)
	if errs != nil {
		if errs.Error() == constants.NoRecordsAffectedError {
			return errors.New(constants.NoRecordsAffectedError)
		}
		return errors.New(constants.AuthenticationFailedError)
	}

	return nil
}
