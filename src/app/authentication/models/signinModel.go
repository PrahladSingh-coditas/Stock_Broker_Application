package models


type BFFSignInUserRequest struct {
	Username string `json:"username" example:"Dharmesh" validate:"required,min=5,max=32"`
	Password string `json:"password" example:"Dk@12345678" validate:"required,strongPassword,min=8,max=20"`
} 

type BFFSignInUserResponse struct {
	Message string `json:"message" example:"user created successfully"`
}

