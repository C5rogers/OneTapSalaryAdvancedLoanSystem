package jwt_auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAuthAccessToken(privateKey string, claims map[string]interface{}) (string, error) {

	mapClaims := jwt.MapClaims(claims)

	_privateKey, err := LoadECDSAPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to load private key: %w", err)
	}
	accessToken, err := sign(mapClaims, _privateKey, "ES256")
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return accessToken, nil
}
