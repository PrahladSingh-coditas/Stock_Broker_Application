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

func NewValidateUserOtpServiceForTest(mockRepo repository.ValidateUserOtpRepository, db *gorm.DB) *ValidateUserOtpService {
	return &ValidateUserOtpService{
		repository: mockRepo,
		db:         db,
	}
}

// this function takes userRequest, fetches the user from db(via repository), performs all otp validations and returns error/ nil
func (service *ValidateUserOtpService) ValidateUserOtp(ctx context.Context, spanCtx context.Context, bffValidateUserOtpRequest models.BFFValidateUserOtpRequest) (string, error) {
	// postgresClinet := utils.GetPostgresClient().GormDB
	userFromDB, err := service.repository.GetUserByUsername(spanCtx, service.db, bffValidateUserOtpRequest.Username)
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

	accessToken, _, err := utils.GenerateToken(bffValidateUserOtpRequest.Username)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
