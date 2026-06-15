package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID           uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"column:username;uniqueIndex" json:"username"`
	Password     string `gorm:"column:password" json:"password"`
	PanCard      string `gorm:"column:panCard;uniqueIndex" json:"panCard"`
	PhoneNumber  uint64 `gorm:"column:phoneNumber" json:"phoneNumber"`
	Email        string `gorm:"column:email;uniqueIndex" json:"email"`
	OtpSent      uint64 `json:"otpSent" gorm:"column:otpSent;default:null"`
	OtpExpiresAt uint64 `json:"otpExpiresAt" gorm:"column:otpExpiresAt;default:null"`
}

type DatabaseConfiguration struct {
	GormDB *gorm.DB
}
