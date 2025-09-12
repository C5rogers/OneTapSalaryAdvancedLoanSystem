package tests

import (
	"testing"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/internal/password"
)

func TestPasswordHash(t *testing.T) {
	hash, err := password.HashPassword("test123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if hash == "" {
		t.Error("Expected hash to be non-empty")
	}
}
