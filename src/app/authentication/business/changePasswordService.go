package business

import (
	"authentication/commons"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"

	"gorm.io/gorm"
)

type ChangePasswordService struct {
	repository repository.ChangePasswordRepository
	db         *gorm.DB //better for testing
}

func NewChangePasswordService(repo repository.ChangePasswordRepository, db *gorm.DB) *ChangePasswordService {
	return &ChangePasswordService{
		repository: repo,
		db:         db,
	}
}

// takes username, newpassword and confirmpassword - compares both passwords and updates in the table if user found and both passwords pass validations
func (service ChangePasswordService) ChangePassword(ctx context.Context, username string, newPassword string, confirmPassword string) error {
	// postgresClient := utils.GetPostgresClient().GormDB //if did this way, testing would be difficult

	_, err := service.repository.GetUserByUsername(ctx, service.db, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return commons.UserNotFoundError
		}
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return commons.HashnigPasswordError
	}

	if err := service.repository.UpdatePassword(ctx, service.db, username, hashedPassword); err != nil {
		return err
	}

	return nil
}
