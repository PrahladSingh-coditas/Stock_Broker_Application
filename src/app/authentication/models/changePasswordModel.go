package models

//just use eqfied ijn confirm poassword
type BFFChangePasswordRequest struct {
	NewPassword     string `json:"newpassword" example:"Sakshi@123" validate:"required,min=8,strongPassword,max=20"`
	ConfirmPassword string `json:"confirmpassword" example:"Sakshi@123" validate:"eqfield=NewPassword"`
}

type BFFChangePasswordResponse struct {
	Message string `json:"message" example:"password updated successfully"`
}
