package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"
	genConstants "stock_broker_application/src/constants"
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

// this function takes userRequest, fetches the user from db(via repository), performs all otp validations and returns error/ nil
func (service *ValidateUserOtpService) ValidateUserOtp(ctx context.Context, spanCtx context.Context, bffValidateUserOtpRequest models.BFFValidateUserOtpRequest) (string, error) {

	postgresClinet := utils.GetPostgresClient()
	tx := postgresClinet.GormDB

	userData, err := service.repository.GetUserByUsername(spanCtx, tx, bffValidateUserOtpRequest.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", constants.UserNotFoundError
		}
		return "", err
	}

	if !utils.CompareUserRequestOTP(userData.OtpSent, bffValidateUserOtpRequest.Otp) {
		return "", constants.IncorrectOTPError
	}

	if !utils.CheckOtpExpiry(userData.OtpExpiresAt, time.Now()) {
		return "", constants.OtpExpiredError
	}

	accessToken, err := utils.GenerateToken(bffValidateUserOtpRequest.Username, constants.PasswordResetPurpose)

	if err != nil {
		return "", constants.TokenCreationFailedError
	}

	return accessToken, nil
}
