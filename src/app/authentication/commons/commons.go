package commons

import (
	"authentication/commons/constants"
	"errors"
)

// Add your common functionalities here.

var UserNotFoundError = errors.New(constants.UserNotFoundError)
var IncorrectPasswordError = errors.New(constants.IncorrectOTPError)
var IncorrectOTPError = errors.New(constants.OtpMismatchError)
var OtpExpiredError = errors.New(constants.OtpExpiredError)
var TokenGenerationError = errors.New(constants.TokenGenerationError)

// constants for returning keys
const (
	Username = "username"
	Password = "password"
	Otp      = "OTP"
	Token = "token"
)
