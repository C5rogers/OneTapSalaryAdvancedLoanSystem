package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Audience    string                 `json:"aud"`
	Issuer      string                 `json:"iss"`
	Expiration  int64                  `json:"exp"`
	Subject     string                 `json:"sub"`
	Email       string                 `json:"email"`
	FullName    string                 `json:"fullName"`
	PhoneNumber string                 `json:"phone_number"`
	IssuedAt    int64                  `json:"iat"`
	Metadata    map[string]interface{} `json:"metadata"`
	jwt.RegisteredClaims
}

func (u *UserClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(u.Expiration, 0)), nil
}

func (u *UserClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(u.IssuedAt, 0)), nil
}

func (u *UserClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (u *UserClaims) GetIssuer() (string, error) {
	return u.Issuer, nil
}

func (u *UserClaims) GetSubject() (string, error) {
	return u.Subject, nil
}

func (u *UserClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{u.Audience}, nil
}
