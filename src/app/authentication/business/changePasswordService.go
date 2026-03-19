package business

import (
	"authentication/models"
	"authentication/repository"
	"stock_broker_application/src/utils"
)

type ChangePasswordService struct {
	repository repository.ChangePasswordRepository
}

func NewChangePasswordService(repository repository.ChangePasswordRepository) *ChangePasswordService {
	return &ChangePasswordService{
		repository: repository,
	}
}

func (service *ChangePasswordService) UpdatePassword(username string, req models.BFFResetPasswordRequest) error {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	return service.repository.UpdatePassword(username, hashedPassword)

}
