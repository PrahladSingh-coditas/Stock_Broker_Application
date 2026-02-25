package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"
)

type ForgotPasswordService struct {
	forgotPasswordRepository repository.ForgotPasswordRepository // we have taken a field that is an interface in repository
}

func NewForgotPasswordService(forgotPasswordRepository repository.ForgotPasswordRepository) *ForgotPasswordService {
	return &ForgotPasswordService{
		forgotPasswordRepository: forgotPasswordRepository, //constructor
	}
}

func (service *ForgotPasswordService) ForgotPassword(ctx context.Context, spanCtx context.Context, bffForgotPasswordRequest models.BFFForgotPasswordRequest) error {
	postgresClinet := utils.GetPostgresClient()
	client := postgresClinet.GormDB

	user, err := service.forgotPasswordRepository.ForgotPassword(spanCtx, client, bffForgotPasswordRequest.Username)
	if err != nil {
		return err
	}

	if user.PanCard != bffForgotPasswordRequest.PanCard || user.PhoneNumber != bffForgotPasswordRequest.PhoneNumber {
		return errors.New(constants.AuthenticationFailedError)
	}

	errs := service.forgotPasswordRepository.GenerateOTP(spanCtx, client, bffForgotPasswordRequest.Username)
	if errs != nil {
		return errs
	}
	return nil
}
