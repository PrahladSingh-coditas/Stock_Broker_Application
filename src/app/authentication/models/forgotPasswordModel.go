package models

import "time"

type BFFForgotPasswordRequest struct {
	Username    string `json:"username" example:"Sanjana" validate:"required,min=5,max=32"`
	PanCard     string `json:"panCard" example:"ABCDE1234F" validate:"required,panCard"`
	PhoneNumber uint64 `json:"phoneNumber" example:"8432805566" validate:"required,min=1000000000,max=9999999999"`
}

type BFFForgotPasswordResponse struct {
	Message      string    `json:"message" example:"user signed in successfully"`
	OtpSent      uint64    `json:"otpSent" gorm:"column:otpSent;default:null"`
	OtpExpiresAt time.Time `json:"otpExpiresAt" gorm:"column:otpExpiresAt;default:null"`
}
