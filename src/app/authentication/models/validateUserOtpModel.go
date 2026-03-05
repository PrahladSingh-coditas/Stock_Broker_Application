package models

type BFFValidateUserOtpRequest struct {
	Username string `json:"username" example:"Dharmesh" validate:"required,min=5,max=32"`
	Otp      string `json:"otp" example:"1234" validate:"required,otp"`
}

type BFFValidateUserOtpResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"token"`
}
