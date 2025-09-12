package tests

import (
	"testing"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
)

func TestErrorResponse(t *testing.T) {
	err := utils.CoalesceError()
	if err == nil {
		t.Error("Expected error to be created")
	}
}
