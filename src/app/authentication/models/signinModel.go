package models

type BFFSigninUserRequest struct {
	Username string `json:"username" example:"Arijit" validate:"required,min=5,max=32"`
	Password string `json:"password" example:"Secure@123" validate:"required,min=8,strongPassword,max=20"`
}

type BFFSigninUserResponse struct {
	Message string `json:"message" example:"user signed in successfully"`
}
