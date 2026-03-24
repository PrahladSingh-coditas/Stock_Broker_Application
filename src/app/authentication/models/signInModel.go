package models

type BFFSigninUserRequest struct {
	Username string `json:"username" example:"Sanjana" validate:"required,min=5,max=32"`
	Password string `json:"password" example:"Sanjana@123" validate:"required,min=8,strongPassword,max=20"`
}

type BFFSigninUserResponse struct {
	Message string `json:"message" example:"user signed in successfully"`
}
