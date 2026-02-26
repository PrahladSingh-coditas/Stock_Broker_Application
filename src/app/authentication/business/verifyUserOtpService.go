package business

import (
	"authentication/commons"
	"authentication/models"
	"authentication/repository"
	"context"
	"stock_broker_application/src/utils"

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

func (service *ValidateUserOtpService) ValidateUserOtp(ctx context.Context, spanCtx context.Context, bffValidateUserOtpRequest models.BFFValidateUserOtpRequest) error {

	userFromDB, errGettingUserFromDB := service.repository.GetUserByUsername(spanCtx, service.db, bffValidateUserOtpRequest.Username)
	if errGettingUserFromDB != nil {
		return commons.UserNotFoundError
	}

	if !utils.CompareUserRequestOTP(userFromDB.OtpSent, bffValidateUserOtpRequest.Otp) {
		return commons.IncorrectOTPError
	}

	if !utils.CheckOtpExpiry(userFromDB.OtpExpiresAt, userFromDB.OtpSent) {
		return commons.OtpExpiredError
	}
	return nil
}
