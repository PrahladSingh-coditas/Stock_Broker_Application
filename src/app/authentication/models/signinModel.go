package models

type BFFSigninUserRequest struct {
	Username string `json:"username" example:"Sakshi" validate:"required,min=5,max=32"`
	Password string `json:"password" example:"Sakshi@123" validate:"required,min=8,strongPassword,max=20"`
}

type BFFSigninUserResponse struct {
	Message      string `json:"message" example:"user logged in successfully"`
	OtpSent      string `json:"otpsent" example:"otp generated successfully"`
	OtpExpiresAt string `json:"otpexpiresat" example:"otp expires in 2 minutes"`
}
