package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, message, code string, status int) error {
	w.WriteHeader(status)

	response := map[string]string{
		"error": message,
		"code":  code,
	}

	return CoalesceError(json.NewEncoder(w).Encode(response), fmt.Errorf("API error: %s", message))
}
