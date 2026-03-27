package commons

import (
	"errors"
	"watchlist/commons/constants"
)

// Add your common functionalities here.

var UserNotFoundError = errors.New(constants.ErrUserNotFound)
var IncorrectPasswordError = errors.New(constants.ErrIncorrectPassword)
var IncorrectOTPError = errors.New(constants.ErrOtpsMismatch)
var OtpExpiredError = errors.New(constants.ErrExpiredOtp)

var TokenGenerationFailed = errors.New(constants.ErrTokenGenerationFailed)

// constants for returning keys
const (
	Username = "username"
	Password = "password"
	Otp      = "OTP"
	Token    = "token"
	Subject  = "sub"
	ScripId = "scripid"
	Watchlist = "watchlist"
	Action = "action"
)
