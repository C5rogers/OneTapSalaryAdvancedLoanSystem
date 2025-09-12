package tests

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env") // adjust path as needed
	os.Exit(m.Run())
}
