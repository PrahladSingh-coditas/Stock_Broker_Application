package constants

// Database Transaction & Connection Errors
const (
	ErrBeginTx            = "failed to begin database transaction: %w"
	ErrCommitTx           = "failed to commit database transaction: %w"
	ErrDBConnectionFailed = "Error connecting to database: %s"
	ErrInternalServer     = "internal server error"
)

// Database Initialization & Config Errors
const (
	ErrDBInitFailed             = "Error initializing database: %s"
	ErrDBMigrationFailed        = "Error migrating database: %s"
	ErrLoadConfigFailed         = "failed to load config: %v"
	ErrPostgresConnectionFailed = "failed to connect with the postgres: %s"
	ErrReadConfigFailed         = "failed to read the config file: %s"
	ErrUnmarshallConfigFailed   = "failed to unmarshal the config file %s"
	ErrJWTConfigReadFailed      = "failed to read the JWT config file %s"
	ErrRedisInitFailed          = "failed to initialize redis %s"
)

const (
	ErrPasswordMinLength    = "Password must be at least 8 characters long."
	ErrPasswordLowercase    = "Password must contain at least one lower letter."
	ErrPasswordUppercase    = "Password must contain at least one uppercase letter."
	ErrPasswordDigit        = "Password must contain at least one digit."
	ErrPasswordSpecialChar  = "Password must contain at least one special character (@$!%*?&)."
	ErrConfirmPasswordMatch = "ConfirmPassword must match Password."
	ErrNewPasswordMatch     = "ConfirmPasword must match NewPassword"
)

const (
	ErrInvalidValue       = "invalid value for %s"
	ErrInvalidPanCard     = "Invalid PAN card format. It should be 5 uppercase letters, followed by 4 digits, and 1 uppercase letter."
	ErrInvalidPhoneNumber = "Phone number must be exactly 10 digits long and contain only numbers."
	ErrFieldRequired      = "%s is required."
	ErrInvalidEmail       = "Invalid value for Email"
)

const (
	ErrHashingPassword      = "error hashing password: %w"
	ErrAuthenticationFailed = "authentication failed"
	ErrUserNotFound         = "user not found"
)

const (
	InvalidTokenError            = "invalid token"
	AuthHeaderMissingError       = "missing authorization header"
	InvalidAuthFormatError       = "invalid authorization format"
	TokenAlreadyBlacklistedError = "token is already blacklisted"
	UnexpectedSigningMethod      = "unexpected signing method"
	InvalidExpClaimsError        = "invalid exp claims"
	InvalidSubClaimsError        = "invalid sub claims"
	TokenExpiredError            = "token is already expired"
	RedisClientNotInitialized    = "redis client not initialized"
)
