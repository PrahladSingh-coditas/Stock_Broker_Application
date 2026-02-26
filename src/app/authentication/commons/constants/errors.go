package constants

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
	ErrUsernameExists    = "username already exists"
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
)

const (
	ErrSignInFailed      = "failed to sign in user"
	ErrUserNotFound      = "user not found"
	ErrIncorrectPassword = "entered password is not correct"
	ErrOtpsMismatch      = "OTPs did not match"
	ErrExpiredOtp        = "OTP expired"
	ErrIncorrectOtp     = "entered OTP is not correct"
)
