package constants

const (
	ErrHeaderMissing         = "Authorization header is missing"
	ErrUsernameNotFoundInJWT = "Username not found in JWT"
)

// Signin and Token generation Errors
const (
	ErrInvalidEmailorPassword = "invalid email or password"
	ErrPasswordMismatch       = "password does not match %w"
	ErrAuthenticationFailed   = "authentication failed"
	ErrTokenGenerationFailed  = "failed to generate authentication tokens %s"
	ErrUserNotFound           = "user not found %w"
	ErrUserNotFoundMsg        = "user not found"
	ErrDatabaseQueryError     = "database query error %w"
	ErrDatabaseQueryErrorMsg  = "database query error"
	ErrRequestFailed          = "request failed"
	ErrPasswordNotMatch       = "password not matched"
	ErrBindingFailed          = "json to struct binding failed"
	ErrValidationFailed       = "validation failed"
	ErrScripNotFoundMsg       = "scrip Id not found"
	ErrNoWatchlistForScripMsg = "no valid watchlists found for user"
	ErrUserAlreadyExistsMsg   = "user already exists"
	ErrFailedToEncrypt        = "failed to encrypt password"
	ErrServer                 = "server"
	InvalidWatchlistIdsError  = "some watchlist ids were invalid"
)

// Request Validation Errors
const (
	ErrInvalidPayload  = "invalid required payload"
	ErrUnexpectedValue = "unexpected value for the field."
)

const (
	ErrInvalidValue  = "invalid value for %s"
	ErrFieldRequired = "%s is required."
)

// Database Transaction & Connection Errors
const (
	ErrBeginTx             = "failed to begin database transaction: %w"
	ErrCommitTx            = "failed to commit database transaction: %w"
	ErrDBConnectionFailed  = "Error connecting to database: %s"
	ErrInternalServer      = "internal server error"
	ErrJWTConfigReadFailed = "failed to read the JWT config file %s"
	OperationFailedError   = "operation failed"
	InvalidTokenError      = "invalid token"
)

// Redis related errors
const (
	EmptyWatchlistIdsError  = "watchlist Ids can't be empty for this operation"
	RedisConnectionError    = "redis connection error"
	RedisUnmarshallingError = "redis: data unmarshalling error"
	RedismarshallingError   = "redis: data marshalling error"
	RedisDataAdditionError  = "redis: data addition error"
	RedisDataDeletionError  = "redis: data deletion error"
)
