package models

type BFFChangePasswordRequest struct {
	NewPassword     string `json:"new_password" validate:"required,min=8,strongPassword,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,eqfield=NewPassword"`
}

type BFFChangePasswordResponse struct {
	Message string `json:"message"`
}

// need to handle the SECFRET KEY in helpers.go/ExtractTokenClaims()
