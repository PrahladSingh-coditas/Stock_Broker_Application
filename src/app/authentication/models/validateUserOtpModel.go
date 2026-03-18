package models

type BFFValidateUserOtpRequest struct {
	Username string `json:"username" validate:"required,min=5,max=32"`
	Otp      uint64 `json:"otp"  validate:"lte=9999,gte=0000"`
}

type BFFValidateUserOtpResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
