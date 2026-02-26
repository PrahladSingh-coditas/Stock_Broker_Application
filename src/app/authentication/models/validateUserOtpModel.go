package models

type BFFValidateUserOtpRequest struct {
	Username string `json:"username" example:"Arijit" validate:"required,min=5,max=32"`
	Otp      string `json:"otp" example:"1234" validate:"required,len=4,numeric"`
}

type BFFValidateUserOtpResponse struct {
	Message string `json:"message"`
}
