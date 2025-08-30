package models

type User struct {
	ID          string `json:"id" graphql:"id"`
	Email       string `json:"email" graphql:"email"`
	FullName    string `json:"full_name" graphql:"full_name"`
	PhoneNumber string `json:"phone_number" graphql:"phone_number"`
	Password    string `json:"password" graphql:"password"`
	Roles       []Role `json:"roles" graphql:"roles"`
}

type RegisteringUser struct {
	FullName       string `json:"full_name" graphql:"full_name" validate:"required"`
	Email          string `json:"email" graphql:"email" validate:"required,email"`
	PhoneNumber    string `json:"phone_number" graphql:"phone_number" validate:"required"`
	Password       string `json:"password" graphql:"password" validate:"required"`
	ProfilePicture string `json:"profile_picture" graphql:"profile_picture"`
	Role           string `json:"role" graphql:"role"`
}
