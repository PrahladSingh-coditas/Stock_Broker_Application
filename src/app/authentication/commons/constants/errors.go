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
	ErrPasswordMismatch       = "password does not match "
	ErrAuthenticationFailed   = "authentication failed"
	ErrTokenGenerationFailed  = "failed to generate authentication tokens "
	ErrInternalServer = "internal server error"
	ErrUserNotFound = "user not found"
)

const(
	ErrBinding = "failed to bind"
)

//otp generation error
const(
	ErrOtpFailed = "failed to generate otp"
)


//otp validation errors
const (
	ErrSignInFailed      = "failed to sign in user" 
	ErrIncorrectPassword = "entered password is not correct"
	ErrOtpsMismatch      = "OTPs did not match"
	ErrExpiredOtp        = "OTP expired"
	ErrIncorrectOtp      = "entered OTP is not correct"
	ErrInvalidOtp        = "OTP must be a 4 digit number"
)

//token generation errors
