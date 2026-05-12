package business

import (
	"authentication/models"
	"authentication/repository"
	"context"
	"fmt"
)

// struct declaration
type CreateUserService struct {
	createUserRepository repository.CreateUserRepository
}

// struct initialisation
func NewCreateUserService(createUserRepository repository.CreateUserRepository) *CreateUserService {
	return &CreateUserService{
		createUserRepository: createUserRepository, // service local var = router local var
	}
}

func (service *CreateUserService) CreateNewUser(ctx context.Context, spanCtx context.Context, bffCreateUserRequest models.BFFCreateUserRequest) error {

	err := service.createUserRepository.CreateNewUser(spanCtx, bffCreateUserRequest)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil

}
