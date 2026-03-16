package models

type BFFChangePasswordRequest struct {
	NewPassword     string `json:"newpassword" example:"SanjanaS@123" validate:"required,min=8,strongPassword,max=20"`
	ConfirmPassword string `json:"confirmPassword" gorm:"column:confirmPassword" example:"SanjanaS@123" validate:"required,min=8,eqfield=NewPassword"`
}

type BFFChangePasswordResponse struct {
	Message string `json:"message" example:"password changed successfully"`
}
