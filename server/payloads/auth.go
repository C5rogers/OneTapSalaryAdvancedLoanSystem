package payloads

import "github.com/go-playground/validator/v10"

type RegisterPayload struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6,max=100"`
	FullName    string `json:"full_name" validate:"required,min=2,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=15"`
	Role        string `json:"role" validate:"required,oneof=user admin"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AccessTokenOutput struct {
	AccessToken string `json:"access_token"`
}

type RegisterOutput struct {
	Message string `json:"message"`
}

func ValidateRegisterPayload(payload *RegisterPayload) error {
	validate := validator.New()
	return validate.Struct(payload)
}

func ValidateLoginPayload(payload *LoginPayload) error {
	validate := validator.New()
	return validate.Struct(payload)
}
