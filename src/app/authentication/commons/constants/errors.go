package constants

// Database Constraint & Index Names
const (
	ErrUniqueConstraintViolation = "duplicate key value violates unique constraint"
	IndexUsersPanCard            = "idx_users_pan_card"
	IndexUsersEmail              = "idx_users_email"
)

// Field Names (JSON/DB)
const (
	FieldPanCard     = "panCard"
	FieldPhoneNumber = "phoneNumber"
	FieldEmail       = "email"
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
	UserNotFoundError            = "user not found"               // try to find by username
	InvalidUsernamePasswordError = "invalid username or password" //invalid userame or password
	PasswordMismatchError        = "password does not match %w"
	AuthenticationFailedError    = "authentication failed"
	TokenGenerationFailedError   = "failed to generate authentication tokens %s"
	MissingCredentialsError      = "username and password are required"
)

// Forgot Password Errors
const (
	InvalidCredentialsError = "invalid pancard and phonenumber"
	DataNotFoundError       = "no user match the credentials"
	NoRecordsAffectedError  = "query had no effect on rows"
	GenerateOtpError        = "error while generating otp"
	DatabaseError           = "database error"
)

//OTP Validation Errors
const (
	IncorrectOTPError    = "Incorrect OTP"
	OtpExpiredError      = "OTP Expired"
	OtpMismatchError     = "OTP Mismatched"
	SigninFailedError    = "Signin Failed Using OTP"
	TokenGenerationError = "Error in Generating Token"
)
