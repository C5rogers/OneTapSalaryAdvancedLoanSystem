package jwt_auth

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func sign(mapClaims jwt.MapClaims, key interface{}, method string) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod(method), mapClaims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseEd25519PrivateKey(pemKey []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}

	edKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("parsed key is not an Ed25519 private key")
	}

	return edKey, nil
}

func LoadECDSAPrivateKey(keyContent string) (*ecdsa.PrivateKey, error) {

	privateKeyStr := strings.ReplaceAll(keyContent, "\\n", "\n")
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing ECDSA private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func LoadECDSAPublicKey(keyContent string) (*ecdsa.PublicKey, error) {
	publicKeyStr := strings.ReplaceAll(keyContent, "\\n", "\n")
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing ECDSA public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid ECDSA public key format")
	}
	return ecdsaPubKey, nil
}
