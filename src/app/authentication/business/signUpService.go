package business

import (
	"authentication/models"
	"authentication/repository"
	"context"
	"fmt"
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

	err := service.createUserRepository.CreateNewUser(spanCtx, bffCreateUserRequest)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil

}
