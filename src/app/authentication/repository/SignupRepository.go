package repository

import (
	"authentication/commons/constants"
	"authentication/models"
	"context"
	"errors"
	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils"
	"strings"
	"time"

	"github.com/pingcap/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CreateUserRepository interface { 
	CreateNewUser(ctx context.Context, db *gorm.DB, bffCreateUserRequest models.BFFCreateUserRequest) error
}

type createUserRepository struct{}

func NewCreateUserRepository() *createUserRepository { // constructor mainly for dependecy injection
	return &createUserRepository{}
}

func (user *createUserRepository) CreateNewUser(ctx context.Context, db *gorm.DB, bffCreateUserRequest models.BFFCreateUserRequest) error {

	start := time.Now()
	logger := logrus.New()
	hashPassword, err := utils.HashPassword(bffCreateUserRequest.Password)
	if err != nil {
		log.Info(constants.ErrFailedToEncrypt)
	}

	bffCreateUserRequest.Password = hashPassword

	NewUser := genericModels.User{
		Username:    bffCreateUserRequest.Username,
		Password:    bffCreateUserRequest.Password,
		PanCard:     bffCreateUserRequest.PanCard,
		PhoneNumber: bffCreateUserRequest.PhoneNumber,
		Email:       bffCreateUserRequest.Email,
	}


	//INSERT INTO users (username, password, email) VALUES ('Arijit', 'hashedpassword', 'abc@gmail.com');
	result := db.WithContext(ctx).Table(constants.UsersTableName).Create(&NewUser)
	if result.Error != nil {
		errorMsgs := result.Error.Error()
		if strings.Contains(errorMsgs, constants.ErrUniqueConstraintViolation) {
			duplicateKeys := []string{}
			if strings.Contains(errorMsgs, constants.IndexUsersPanCard) {
				duplicateKeys = append(duplicateKeys, constants.FieldPanCard)
			}
			if strings.Contains(errorMsgs, constants.IndexUsersEmail) {
				duplicateKeys = append(duplicateKeys, constants.FieldEmail)
			}

			if len(duplicateKeys) > 0 {
				return errors.New(strings.Join(duplicateKeys, ",") + constants.ErrDuplicateEntry)
			}

			return errors.New(constants.ErrUsernameExists)
		}
		return result.Error
	}

	logger.WithFields(logrus.Fields{
		"user":    bffCreateUserRequest.Email,
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserCreationSuccessMsg)

	return nil
}
