package models

type SignUpPayload struct {
	FullName       string `json:"full_name" graphql:"full_name" validate:"required"`
	Email          string `json:"email" graphql:"email" validate:"required,email"`
	PhoneNumber    string `json:"phone_number" graphql:"phone_number" validate:"required"`
	Password       string `json:"password" graphql:"password" validate:"required"`
	ProfilePicture string `json:"profile_picture" graphql:"profile_picture"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// login
type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	Role string `json:"role" graphql:"role"`
	Name string `json:"name" graphql:"name"`
}

type LoginArgs struct {
	Credential LoginInput
}

type LoginPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            LoginArgs              `json:"input"`
}

type LoginPayloadNew struct {
	Email    string `json:"email" graphql:"email"`
	Password string `json:"password" graphql:"password"`
}

type AccessTokenOutput struct {
	AccessToken string `json:"access_token" graphql:"access_token"`
}

type ChapaPaymentPayload struct {
	BookID string `json:"book_id" graphql:"book_id"`
}

type ChapaPaymentResponse struct {
	PaymentURL string `json:"payment_url" graphql:"payment_url"`
	PaymentID  string `json:"payment_id" graphql:"payment_id"`
}
