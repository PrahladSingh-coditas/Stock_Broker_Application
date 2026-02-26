package business

import (
	"authentication/models"
	"authentication/repository"
	"context"
	"fmt"
	genericErrors "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
)

type CreateUserService struct {
	createUserRepository repository.CreateUserRepository
}

func NewCreateUserService(createUserRepository repository.CreateUserRepository) *CreateUserService {
	return &CreateUserService{
		createUserRepository: createUserRepository,
	}

}

func (service *CreateUserService) CreateNewUser(ctx context.Context, spanCtx context.Context, bffCreateUserRequest models.BFFCreateUserRequest) error {
	postgresClinet := utils.GetPostgresClient()
	tx := postgresClinet.GormDB.Begin()

	if tx.Error != nil {
		return fmt.Errorf(genericErrors.ErrBeginTx, tx.Error)
	}

	err := service.createUserRepository.CreateNewUser(spanCtx, tx, bffCreateUserRequest)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf(genericErrors.ErrCommitTx, err)
	}

	return nil

}
