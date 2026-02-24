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
	ErrUserNotFound            = "user not found"               // try to find by username
	ErrInvalidUsernamePassword = "invalid username or password" //invalide u or p
	ErrPasswordMismatch      = "password does not match %w"
	ErrAuthenticationFailed  = "authentication failed"
	ErrTokenGenerationFailed = "failed to generate authentication tokens %s"
	ErrMissingCredentials    = "username and password are required"
)
