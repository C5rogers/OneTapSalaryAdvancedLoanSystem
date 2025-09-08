package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/payloads"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	claims, ok := r.Context().Value("user").(*payloads.UserClaims)
	if !ok {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	userEmail, ok := claims.Metadata["x-auth-user-id"].(string)
	if !ok {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}
	userRole, ok := claims.Metadata["x-auth-role"].(string)
	if !ok {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	registeringUser, err := s.DB.FindUserByEmail(userEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusBadRequest)
		} else {
			return utils.SendErrorResponse(w, "error checking user", "database_error", http.StatusInternalServerError)
		}
	}

	if registeringUser == nil || registeringUser.Role != "admin" || registeringUser.Role != userRole {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	payload := payloads.RegisterPayload{}
	err = json.Unmarshal(reqBody, &payload)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	if err := payloads.ValidateRegisterPayload(&payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.SendErrorResponse(w, validationErrors.Error(), "validation_error", http.StatusBadRequest)
	}

	existingUser, err := s.DB.FindUserByEmail(payload.Email)
	if err != nil && err.Error() != "record not found" {
		return utils.SendErrorResponse(w, "error checking existing user", "database_error", http.StatusInternalServerError)
	}
	if existingUser != nil {
		return utils.SendErrorResponse(w, "email already in use", "email_in_use", http.StatusBadRequest)
	}

	existingPhoneUser, err := s.DB.FindUserByPhoneNumber(payload.PhoneNumber)
	if err != nil && err.Error() != "record not found" {
		return utils.SendErrorResponse(w, "error checking existing user", "database_error", http.StatusInternalServerError)
	}
	if existingPhoneUser != nil {
		return utils.SendErrorResponse(w, "phone number already in use", "phone_number_in_use", http.StatusBadRequest)
	}

	newUser := &models.User{
		Email:       payload.Email,
		FullName:    payload.FullName,
		PhoneNumber: payload.PhoneNumber,
		Role:        payload.Role,
	}
	newUser.SetPassword(payload.Password)

	if err := s.DB.CreateUser(newUser); err != nil {
		return utils.SendErrorResponse(w, "error creating user", "database_error", http.StatusInternalServerError)
	}

	responseData, err := json.Marshal(&payloads.RegisterOutput{Message: "user registered successfully"})
	if err != nil {
		return utils.SendErrorResponse(w, "error creating user", "database_error", http.StatusInternalServerError)
	}

	if _, err := w.Write(responseData); err != nil {
		return err
	}
	return nil
}
