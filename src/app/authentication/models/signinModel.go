package models

type BFFSignInRequest struct {
	Username string `json:"username"  validate:"required,min=5,max=32"`
	Password string `json:"password"  validate:"required,min=8,max=20"`
}

type BFFSignInResponse struct {
	Message string `json:"message" example:"Signed in successfully"`
}
