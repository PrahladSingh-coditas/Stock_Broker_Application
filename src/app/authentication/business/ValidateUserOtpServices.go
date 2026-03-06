package business

import (
	"authentication/commons"
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"
	"time"

	"gorm.io/gorm"
)

type ValidateUserOtpService struct {
	repository repository.ValidateUserOtpRepository
}

func NewValidateUserOtpService(repository repository.ValidateUserOtpRepository) *ValidateUserOtpService {
	return &ValidateUserOtpService{
		repository: repository,
	}
}

// this function takes userRequest, fetches the user from db(via repository), performs all otp validations and returns error/ nil and token
func (service *ValidateUserOtpService) ValidateUserOtp(ctx context.Context, spanCtx context.Context, bffValidateUserOtpRequest models.BFFValidateUserOtpRequest) (string, error) {
	postgresClient := utils.GetPostgresClient()
	client := postgresClient.GormDB

	userFromDB, err := service.repository.GetUserByUsername(spanCtx, client, bffValidateUserOtpRequest.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", commons.UserNotFoundError //errors.New(constants.ErrUserNotFound)
		}
		return "", err
	}

	if !utils.CompareUserRequestOTP(userFromDB.OtpSent, bffValidateUserOtpRequest.Otp) {
		return "", commons.IncorrectOTPError //errors.New(constants.ErrIncorrectOtp)
	}

	if !utils.CheckOtpExpiry(userFromDB.OtpExpiresAt, time.Now()) {
		return "", commons.OtpExpiredError //errors.New(constants.ErrExpiredOtp)
	}

	tokenString, err := utils.GenerateToken(bffValidateUserOtpRequest.Username, constants.Purpose)
	if err != nil {
		return "", commons.TokenGenerationFailed
	}

	return tokenString, nil
}
