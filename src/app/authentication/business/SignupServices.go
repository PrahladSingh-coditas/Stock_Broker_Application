package business

import (
	"authentication/models"
	"authentication/repository"
	"context"
	"fmt"
	// genericErrors "stock_broker_application/src/constants"
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
	postgresClient := utils.GetPostgresClient()
	client:= postgresClient.GormDB


	err := service.createUserRepository.CreateNewUser(spanCtx, client, bffCreateUserRequest)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil

}
