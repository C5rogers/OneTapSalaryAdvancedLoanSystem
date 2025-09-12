package tests

import (
	"os"
	"testing"

	jwt_auth "github.com/c5rogers/one-tap/salary-advance-loan-system/internal/jwt-auth"
)

func TestJWTAuthInit(t *testing.T) {
	// Attempt to load the ECDSA public key
	// This will also initialize the JWT authentication system
	jwtPublickKeyPayth, err := os.ReadFile(os.Getenv("CONFIG_SERVER__JWT_PUBLIC_KEY_PATH"))
	if err != nil {
		t.Error("Expected no error reading JWT public key file, got", err)
	}

	_, err = jwt_auth.LoadECDSAPublicKey(string(jwtPublickKeyPayth))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
