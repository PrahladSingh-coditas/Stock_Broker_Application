package models

type BFFChangePasswordRequest struct {
	NewPassword     string `json:"password" example:"Xyz@12345" validate:"required,strongPassword,max=20,min=8"`
	ConfirmPassword string `json:"confirmPassword" example:"Xyz@12345" gorm:"column:confirmPassword" validate:"required,min=8,eqfield=NewPassword"`
}

type BFFChangePasswordResponse struct {
	Message string `json:"message" example:"Password changed successfully"`
}
