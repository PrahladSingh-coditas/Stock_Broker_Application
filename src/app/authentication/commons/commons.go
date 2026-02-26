package commons

import (
	"authentication/commons/constants"
	"errors"
)

// Add your common functionalities here.
var UserNotFoundError = errors.New(constants.ErrUserNotFound)
var IncorrectPasswordError = errors.New(constants.ErrIncorrectPassword)
var IncorrectOTPError = errors.New(constants.ErrOtpsMismatch)
var OtpExpiredError = errors.New(constants.ErrExpiredOtp)

const (
	Username = "username"
	Password = "password"
	Otp      = "OTP"
)
