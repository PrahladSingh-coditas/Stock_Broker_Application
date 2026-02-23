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

// Success message for user
const (
	UserCreationSuccessMsg = "User created successfully"
	UserLoggedInSuccessMsg = "User logged in successfully"
)

// error
const (
	ErrInvalidCredentials = " user not signup "
	ErrUsernameNotFound   = " username Not Found"
	ErrLoginFailed        = "Login Failed"
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
