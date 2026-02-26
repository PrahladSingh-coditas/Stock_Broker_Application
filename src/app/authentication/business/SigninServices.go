package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repository"
	"context"
	"errors"
	"stock_broker_application/src/utils"
)

type SignInService struct {
	signinRepository repository.SignInRepository
}

func NewSignInService(signinRepository repository.SignInRepository) *SignInService {
	return &SignInService{
		signinRepository: signinRepository,
	}
}

func (service *SignInService) SignIn(ctx context.Context, spanCtx context.Context, bffSignInRequest models.BFFSignInRequest) error {
	postgresClinet := utils.GetPostgresClient().GormDB
	// tx := postgresClinet.GormDB.Begin()

	userFromDB, errGetUserFromDB := service.signinRepository.GetUserByUsername(spanCtx, postgresClinet, bffSignInRequest.Username)
	if errGetUserFromDB != nil {
		return errors.New(constants.ErrUserNotFound)
	}

	//mock otp testing by updating the otpSent and otpExpiresAt fields in db table when signed in for validation task
	//added this fort testing purpose, to manually update the request time and a fixed otp
	// otp := uint64(1010)
	// expiry := uint64(time.Now().Add(2 * time.Minute).Unix())

	// postgresClinet.Model(&userFromDB).Updates(map[string]interface{}{
	// 	"otpSent":      otp,
	// 	"otpExpiresAt": expiry,
	// })
	// fmt.Println("mock otp generated: ", otp)
	// fmt.Println("otp expiry epoch time: ", expiry)

	if !utils.CompareHashPassword(userFromDB.Password, bffSignInRequest.Password) {
		return errors.New(constants.ErrIncorrectPassword)
	}
	return nil
}
