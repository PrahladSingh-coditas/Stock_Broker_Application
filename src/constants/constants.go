package constants

const (
	RunningServerPort = "Running Server on port : %v"
)

const (
	PANCardRegex     = `^[A-Z]{5}[0-9]{4}[A-Z]{1}$`
	PasswordRegex    = `^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
	EmailRegex       = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z]{2,}$`
	UppercaseRegex   = `[A-Z]`
	DigitRegex       = `\d`
	SpecialCharRegex = `[@$!%*?&]`
	LowercaseRegex   = `[a-z]`
	OtpRegexp        = "^[0-9]{4}$"
)

const (
	FieldPassword        = "Password"
	FieldConfirmPassword = "ConfirmPassword"
	FieldPanCard         = "PanCard"
	FieldStrongPassword  = "strongPassword"
	FieldPhoneNumber     = "PhoneNumber"
	FieldEmail           = "Email"
	FieldUsername        = "Username"
	FieldOtp             = "Otp"
)

// Migration success Message
const (
	MsgDBMigrationSuccess = "Database migration completed successfully!"
)

const (
	Postgres = "postgres"
	JWT      = "jwt"
	Yaml     = "yaml"
)

const (
	DSNString = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s"
)

// Database Keys
const (
	Interal = "internal"
)

// Paths
const (
	BaseConfig = "../../config"
	RootConfig = "./src/config"
)

// Origin
const (
	AllowedOrigin = "*"
)

// Method
const (
	POST = "POST"
	GET  = "GET"
)

// Header
const (
	Origin        = "Origin"
	ContentType   = "Content-type"
	Authorization = "Authorization"
)

// redis
const (
	RedisAddr     = "127.0.0.1:6379"
	RedisPassword = ""
	RedisDB       = 0
)
