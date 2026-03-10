package business

import (
	"authentication/commons/constants"
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

// this function takes userRequest, fetches the user from db(via repository), performs all otp validations and returns access token and error/ nil
func (service *ValidateUserOtpService) ValidateUserOtp(ctx context.Context, spanCtx context.Context, bffValidateUserOtpRequest models.BFFValidateUserOtpRequest) (string, error) {
	postgresClinet := utils.GetPostgresClient()
	client := postgresClinet.GormDB
	userFromDB, err := service.repository.GetUserByUsername(spanCtx, client, bffValidateUserOtpRequest.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New(constants.UserNotFoundError)
		}
		return "", err
	}

	otp, err := strconv.ParseUint(bffValidateUserOtpRequest.Otp, 10, 64)
	if err != nil {
		return "", err
	}
	if otp != userFromDB.OtpSent {
		return "", errors.New(constants.IncorrectOTPError)
	}

	if userFromDB.OtpExpiresAt > uint64(time.Now().Unix()) {
		return "", errors.New(constants.OtpExpiredError)
	}

	token, _, err := utils.GenerateToken(bffValidateUserOtpRequest.Username)
	if err != nil {
		return "", errors.New(constants.TokenGenerationError)
	}

	return token, nil
}
