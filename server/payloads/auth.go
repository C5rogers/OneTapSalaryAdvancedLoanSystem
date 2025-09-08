package payloads

import (
	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
	"github.com/go-playground/validator/v10"
)

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

type Result struct {
	Customer models.Customer `json:"customer"`
	Rating   RatingBreakdown `json:"rating"`
}

type ValidateCustomerOutput struct {
	ValidCustomerLogs   []models.ValidationLog `json:"valid_customer_logs"`
	RegisteredCustomers []models.Customer      `json:"registered_customers"`
	InvalidCustomerLogs []models.ValidationLog `json:"invalid_customer_logs"`
}

type ProcessTransactionsOutput struct {
	Ratings []Result `json:"ratings"`
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
