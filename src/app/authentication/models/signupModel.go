package models

type BFFCreateUserRequest struct {
	Username        string `json:"username" example:"Arijit" validate:"required,min=5,max=32"`
	Password        string `json:"password"  example:"Secure@123" validate:"required,min=8,strongPassword,max=20"`
	ConfirmPassword string `json:"confirmPassword" gorm:"column:confirmPassword" example:"Secure@123" validate:"required,min=8,eqfield=Password"`
	PanCard         string `json:"panCard" gorm:"column:panCard" example:"EQZRP1234P" validate:"required,panCard"`
	PhoneNumber     uint64 `json:"phoneNumber" gorm:"column:phoneNumber" example:"7568912340" validate:"required,min=1000000000,max=9999999999"`
	Email           string `json:"email" example:"arijit@gmail.com" validate:"required,Email"`
}

type BFFCreateUserResponse struct {
	Message string `json:"message" example:"user created successfully"`
}
