package constants

//Authentications API URL Keys
const (
	ServiceName       = "authentication"
	PortDefaultValude = 8080
)

// Database table name & field names for users
const (
	UsersTableName = "users"
	Fieldemail     = "email"
)

const(
	UsernameField = "username = ?"
)

// Success message for user login
const (
	UserCreationSuccessMsg = "User created successfully"
	UserLoggedInSuccessMsg = "User logged in successfully"
)

//success message for user otp generation
const (
	UserOtpGeneratedSuccess = "Otp generated successfully."
	UserOtpExpiryMsg = "Otp will expire in 2 minutes."
)

//success message for user otp validation
const OtpValidatedSuccessMsg = "OTP validated successfully"

//success message for password change
const PasswordChangedSuccess = "Password changed successfully"

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

//purpose for JWT
const Purpose = "password_reset"