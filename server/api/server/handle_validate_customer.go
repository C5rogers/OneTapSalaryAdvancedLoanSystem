package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/payloads"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
	"gorm.io/gorm"
)

func (s *Server) HandleValidateCustomers(w http.ResponseWriter, r *http.Request) error {

	claims, ok := r.Context().Value("user").(*payloads.UserClaims)
	if !ok {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	userEmail, ok := claims.Metadata["x-auth-user-id"].(string)
	if !ok {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	registeringUser, err := s.DB.FindUserByEmail(userEmail)
	if err != nil {
		if err.Error() == "record not found" {
			return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusBadRequest)
		}
		return utils.SendErrorResponse(w, "error checking user", "database_error", http.StatusInternalServerError)
	}

	if registeringUser == nil {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	userRole, ok := claims.Metadata["x-auth-role"].(string)
	if !ok {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	if userRole != registeringUser.Role {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	logs, validCustomers, err := utils.ValidateCustomers("data/customers.json", "sample_customers.csv")
	if err != nil {
		return utils.SendErrorResponse(w, "error validating customers: "+err.Error(), "validation_error", http.StatusInternalServerError)
	}

	// save only valid ones into DB
	for _, c := range validCustomers {

		_, err := s.DB.FindCustomerByAccountNumber(c.AccountNumber)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// create new customer
				if err := s.DB.CreateCustomer(&c); err != nil {
					return utils.SendErrorResponse(w, "error saving customer", "database_error", http.StatusInternalServerError)
				}
				continue
			}
			return utils.SendErrorResponse(w, "error checking existing customer", "database_error", http.StatusInternalServerError)
		}
	}

	// separate the valid customer log and invalid customer logs
	var invalidCustomerLogs []models.ValidationLog
	for _, log := range logs {
		if !log.Verified {
			invalidCustomerLogs = append(invalidCustomerLogs, log)
		}
	}

	var validLogs []models.ValidationLog
	for _, log := range logs {
		if log.Verified {
			validLogs = append(validLogs, log)
		}
	}

	// record the log with the current timestamp to file
	logFileName := "logs/validation_log_" + time.Now().Format("20060102150405") + ".json"
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			return utils.SendErrorResponse(w, "error creating log directory", "file_error", http.StatusInternalServerError)
		}
	}
	logData, _ := json.Marshal(logs)
	if err := os.WriteFile(logFileName, logData, 0644); err != nil {
		return utils.SendErrorResponse(w, "error saving validation log", "file_error", http.StatusInternalServerError)
	}

	responseData, _ := json.Marshal(payloads.ValidateCustomerOutput{
		ValidCustomerLogs:   validLogs,
		RegisteredCustomers: validCustomers,
		InvalidCustomerLogs: invalidCustomerLogs,
	})

	if _, err := w.Write(responseData); err != nil {
		return err
	}
	return nil
}
