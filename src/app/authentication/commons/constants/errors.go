package constants

import "errors"

// Database Constraint & Index Names
const (
	ErrUniqueConstraintViolation = "duplicate key value violates unique constraint"
	IndexUsersPanCard            = "idx_users_pan_card"
	IndexUsersEmail              = "idx_users_email"
)

// Field Names (JSON/DB)
const (
	FieldPanCard = "panCard"
	FieldEmail   = "email"
)

// Duplicate Entry Errors
const (
	ErrDuplicateEntry    = "already exists"
	ErrUsernameExists    = "usernamealready exists"
	ErrUserAlreadyExists = "user already exists"
)

// General Errors
const (
	ErrConflict           = "conflict"
	ErrUserCreationFailed = "failed to create user"
)

// Request Validation Errors
const (
	ErrInvalidPayload  = "invalid required payload"
	ErrUnexpectedValue = "unexpected value for the field."
)

//Encrypt & Decrypt Erros
const (
	ErrFailedToEncrypt = "falied to encrpyt password"
)

//Signin and Token generation Errors
const (
	ErrInvalidEmailorPassword = "invalid email or password"
	ErrPasswordMismatch       = "password does not match %w"
	ErrAuthenticationFailed   = "authentication failed"
	ErrTokenGenerationFailed  = "failed to generate authentication tokens %s"
	ErrUserNotFound           = "user not found %w"
	ErrUserNotFoundMsg        = "user not found"
	ErrDatabaseQueryError     = "database query error %w"
	ErrDatabaseQueryErrorMsg  = "database query error"
	ErrPasswordNotMatch       = "password not matched"
	ErrBindingFailed          = "json to struct binding failed"
	ErrValidationFailed       = "validation failed"
)

//redis and token related errors
const (
	//RedisConnectionError = "error in redis connection"
	RedisSetOperationError = "error in redis set operation"
	RedisExistsOperationError = "error in redis exists operation"
	RedisGetOperationError = "error in redis get operation"
	InvalidTokenError    = "invalid token"
	LogoutFailed         = "logout failed"
	OperationFailed         = "operation failed"
)

//Otp related errors
const (
	ErrSignInFailed            = "failed to sign in user"
	ErrOtpsMismatch            = "OTPs did not match"
	ErrExpiredOtp              = "OTP expired"
	ErrIncorrectOtp            = "entered OTP is not correct"
	ErrInvalidOtp              = "OTP must be a 4 digit number"
	ErrTokenCreationFailed     = "token generation failed"
	ErrHeaderMissing           = "Authorization header is missing"
	ErrUsernameNotFoundInJWT   = "Username not found in JWT"
	ErrExpiryTimeNotFoundInJWT = "expiry time not found in JWT"
	ErrPurposeNotFoundInHeader = "purpose not found inside the request"
	ErrPurposeNotMatched       = "purpose is not matching"
)

var UserNotFoundError = errors.New(ErrUserNotFoundMsg)
var IncorrectPasswordError = errors.New(ErrIncorrectOtp)
var IncorrectOTPError = errors.New(ErrOtpsMismatch)
var OtpExpiredError = errors.New(ErrExpiredOtp)
var TokenCreationFailedError = errors.New(ErrTokenCreationFailed)
var DatabaseQueryError = errors.New(ErrDatabaseQueryErrorMsg)
var PasswordEncryptFailedError = errors.New(ErrFailedToEncrypt)
var UsernameNotFoundInJWTError = errors.New(ErrUsernameNotFoundInJWT)
