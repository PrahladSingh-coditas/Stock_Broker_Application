package constants

//Authentications API URL Keys
const (
	ServiceName       = "authentication"
	PortDefaultValude = 8080
)

//Database table name & field names for users
const (
	UsersTableName = "users"
	Fieldemail     = "email"
	FieldUsername  = "username = ?"
	FieldPassword  = "password"
	Username       = "username"
)

// Success message for user
const (
	UserCreationSuccessMsg              = "User created successfully"
	UserLoggedInSuccessMsg              = "User logged in successfully"
	UserReadSuccessMsg                  = "User read successfully"
	CredentialMatchSuccessMsg           = "Credential matched successfully"
	ForgotPasswordGenerateOtpSuccessMsg = "Forgot password OTP sent successfully, will expire in 2 minutes"
	OtpValidatedSuccessMsg              = "OTP Validation Successful"
	TokenGeneratedSuccessMsg            = "Token Generated Successfully"
	PasswordChangeSuccessMsg            = "User Password Changed Successfully"
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

const (
	Username = "username"
	Password = "password"
	Otp      = "OTP"
)
