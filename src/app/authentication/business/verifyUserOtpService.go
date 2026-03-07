package business

import (
	"authentication/commons"
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
	db         *gorm.DB
}

func NewValidateUserOtpService(repository repository.ValidateUserOtpRepository, db *gorm.DB) *ValidateUserOtpService {
	return &ValidateUserOtpService{
		repository: repository,
		db:         db,
	}
}

func (service *ValidateUserOtpService) ValidateUserOtp(ctx context.Context, spanCtx context.Context, bffValidateUserOtpRequest models.BFFValidateUserOtpRequest) (string, error) {

	userFromDB, err := service.repository.GetUserByUsername(spanCtx, service.db, bffValidateUserOtpRequest.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", commons.UserNotFoundError
		}
		return "", err
	}

	if !utils.CompareUserRequestOTP(userFromDB.OtpSent, bffValidateUserOtpRequest.Otp) {
		return "", commons.IncorrectOTPError
	}

	if !utils.CheckOtpExpiry(userFromDB.OtpExpiresAt, time.Now()) {
		return "", commons.OtpExpiredError
	}

	token, err := utils.GeneratePasswordResetToken(userFromDB.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}
