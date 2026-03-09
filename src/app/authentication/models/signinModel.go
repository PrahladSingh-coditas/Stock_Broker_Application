package models

type BFFSignInRequest struct {
	Username string `json:"username" example:"Arijit" validate:"required,min=5,max=32"`
	Password string `json:"password" example:"Secure@123" validate:"required,min=8,max=20"`
}

type BFFSignInResponse struct {
	Message string `json:"message" example:"Signed in successfully"`
}
