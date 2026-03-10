package business

import (
	"authentication/commons"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"
	"strconv"
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
func (service *ValidateUserOtpService) ValidateUserOtp(ctx context.Context, spanCtx context.Context, bffValidateUserOtpRequest models.BFFValidateUserOtpRequest) error {
	postgresClinet := utils.GetPostgresClient()
	client := postgresClinet.GormDB
	userFromDB, err := service.repository.GetUserByUsername(spanCtx, client, bffValidateUserOtpRequest.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return commons.UserNotFoundError
		}
		return err
	}

	parsedOTP, err := strconv.ParseUint(bffValidateUserOtpRequest.Otp, 10, 64)
	if err != nil {
		return err
	}
	if userFromDB.OtpSent != parsedOTP {
		return commons.IncorrectOTPError
	}

	if !utils.CheckOtpExpiry(userFromDB.OtpExpiresAt, time.Now()) {
		return commons.OtpExpiredError //errors.New(constants.ErrExpiredOtp)
	}
	return nil
}
