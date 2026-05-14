package constants

// Database Transaction & Connection Errors
const (
	ErrBeginTx                 = "failed to begin database transaction: %w"
	ErrCommitTx                = "failed to commit database transaction: %w"
	ErrDBConnectionFailed      = "Error connecting to database: %s"
	ErrInternalServer          = "internal server error"
	RedisConnectionFailedError = "Connection Failed: %s"
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
	RedisConnectionError        = "Error connecting to Redis: %s"
)

const (
	ErrPasswordMinLength    = "Password must be at least 8 characters long."
	ErrPasswordLowercase    = "Password must contain at least one lower letter."
	ErrPasswordUppercase    = "Password must contain at least one uppercase letter."
	ErrPasswordDigit        = "Password must contain at least one digit."
	ErrPasswordSpecialChar  = "Password must contain at least one special character (@$!%*?&)."
	ErrConfirmPasswordMatch = "ConfirmPassword must match Password."
)

const (
	ErrInvalidValue         = "invalid value for %s"
	ErrInvalidPanCard       = "Invalid PAN card format. It should be 5 uppercase letters, followed by 4 digits, and 1 uppercase letter."
	ErrInvalidPhoneNumber   = "Phone number must be exactly 10 digits long and contain only numbers."
	ErrFieldRequired        = "%s is required."
	ErrFieldRequiredIf      = "%s is required when %s is %s."
	ErrFieldExcludedIf      = "%s must be excluded when %s is %s."
	ErrInvalidEmail         = "Invalid value for Email"
	InvalidScripFormatError = "Invalid scripId format (%s). Eg: 'NSE/BSE_2000'"
)
