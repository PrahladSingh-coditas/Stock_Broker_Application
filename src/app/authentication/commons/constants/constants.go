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

// Success message for user
const (
	UserCreationSuccessMsg = "User created successfully"
	UserLoggedInSuccessMsg = "User logged in successfully"
)

//success message for user with otp
const (
	UserLoggedOtpGeneratedSuccess = "User logged in & otp generated successfully"
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
