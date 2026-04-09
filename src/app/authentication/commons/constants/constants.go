package constants

//Authentications API URL Keys
const (
	ServiceName       = "authentication"
	PortDefaultValude = 8080
)

// Database table name & field names for users
const (
	UsersTableName    = "users"
	Fieldemail        = "email"
	UsernameCondition = "username = ? "
	Password          = "password"
	Database          = "database"
)

//Logger related fields
const (
	User       = "username"
	Latency    = "latency"
	Otp        = "OTP"
	Token      = "token"
	Bearer     = "Bearer "
	Subject    = "sub"
	Purpose    = "purpose"
	Expiry     = "exp"
	ExpiryTime = "expiry_time"
	Status     = "status"
	Error      = "error"
	Redis      = "redis"
	Header     = "header"
)

// Success message for user
const (
	UserCreationSuccessMsg = "User created successfully"
	UserLoggedInSuccessMsg = "User logged in successfully"
	UserDataFetchedMsg     = "User data fetched successfully"
	UserOtpFetchedMsg      = "User Otp fetched successfully"
	OtpValidatedSuccessMsg = "OTP validated successfully"
	PasswordUpdateMsg      = "Password updated successfully"
	LogoutSuccessMsg       = "logout successfully"
)

//Otp related messages
const (
	OtpSentAndExpiryMsg   = "OTP sent successfully, OTP will expire in next %d minutes"
	OtpTimeLimitInSeconds = 120
	OtpTimeLimitInMinutes = 2
	OTPSent               = "OtpSent"
	OTPExpiryTime         = "OtpExpiresAt"
)

//Swagger Titile
const SwaggerTitle = "Stock Broker Application API"

const EmailorPasswordField = "email_or_password"

//Cookies
const (
	Name     = "refresh_token"
	Time     = 30 * 24 * 60 * 60
	Path     = "/"
	Domain   = ""
	Secure   = true
	HttpOnly = true
)

//token related messages
const (
	PasswordResetPurpose = "password_reset"
	BlacklistedCacheKey  = "BLACKLISTED_TOKEN_%s"
)
