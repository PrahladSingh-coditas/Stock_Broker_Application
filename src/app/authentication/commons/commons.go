package commons

import (
	"authentication/commons/constants"
	"errors"
	genericConstants "stock_broker_application/src/constants"
)

// Add your common functionalities here.
var UserNotFoundError = errors.New(constants.ErrUserNotFound)
var IncorrectPasswordError = errors.New(constants.ErrIncorrectPassword)
var IncorrectOTPError = errors.New(constants.ErrOtpsMismatch)
var OtpExpiredError = errors.New(constants.ErrExpiredOtp)
var InvalidTokenError = errors.New(constants.ErrInvalidToken)
var ConfirmPasswordMismatchError = errors.New(genericConstants.ErrConfirmPasswordMatch)
var NewPasswordMismatchError = errors.New(genericConstants.ErrNewPasswordMatch)
var HashnigPasswordError = errors.New(genericConstants.ErrHashingPassword)

// constants for returning keys
const (
	Username = "username"
	Password = "password"
	Otp      = "OTP"
	Token    = "token"
)
