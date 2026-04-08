package constants

//Authentications API URL Keys
const (
	ServiceName       = "authentication"
	PortDefaultValude = 8081
)

// Database table name & field names for users
const (
	UsersTableName = "users"
	Fieldemail     = "email"
	Username       = "username = ?"
)

// Success message for user
const (
	UserCreationSuccessMsg    = "User created successfully"
	UserLoggedInSuccessMsg    = "User logged in successfully"
	OtpValidatedSuccessMsg    = "OTP validated successfully"
	PasswordChangedSuccessMsg = "Password changed succefully"
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
