package constants

// Database Constraint & Index Names
const (
	ErrUniqueConstraintViolation = "duplicate key value violates unique constraint"
	IndexUsersPanCard            = "idx_users_pan_card"
	IndexUsersEmail              = "idx_users_email"
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
	ScripIdNotFoundError      = "scripId not found"
	UserNotFoundError         = "user not found" // try to find by username
	AuthenticationFailedError = "authentication failed"
	QueryError                = "error in executing query"
	InvalidActionTypeError    = "invalid action type, must be 'ADD' 'GET' or 'DEL'"
	WatchlistNotFoundError    = "watchlist not found"
)

//error
const (
	ServerError = "Something went wrong"
)
