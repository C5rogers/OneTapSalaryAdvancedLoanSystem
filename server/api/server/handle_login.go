package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
	jwt_auth "github.com/c5rogers/one-tap/salary-advance-loan-system/internal/jwt-auth"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/payloads"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/security"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) error {

	ip := security.GetIP(r.RemoteAddr)
	if !s.RateLimiter.Allow(ip) {
		return utils.SendErrorResponse(w, "too many requests", "rate_limit_exceeded", http.StatusTooManyRequests)
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	payload := payloads.LoginPayload{}

	err = json.Unmarshal(reqBody, &payload)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	if err := payloads.ValidateLoginPayload(&payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.SendErrorResponse(w, validationErrors.Error(), "validation_error", http.StatusBadRequest)
	}

	user, err := s.DB.FindUserByEmail(payload.Email)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(w, "username or password incorrect", "bad_credentials", http.StatusBadRequest)
		} else {
			return utils.SendErrorResponse(w, "error retrieving user", "database_error", http.StatusInternalServerError)
		}
	}

	if user == nil {
		return utils.SendErrorResponse(w, "username or password incorrect", "bad_credentials", http.StatusBadRequest)
	}

	if user == (&models.User{}) {
		return utils.SendErrorResponse(w, "username or password incorrect", "bad_credentials", http.StatusBadRequest)
	}

	if !user.CheckPassword(payload.Password) {
		return utils.SendErrorResponse(w, "username or password incorrect", "bad_credentials", http.StatusBadRequest)
	}

	accessTokenPayload := map[string]interface{}{
		"aud":          user.Email,
		"iss":          "https://one-tap.backend.com",
		"exp":          time.Now().Add(72 * time.Hour).Unix(),
		"sub":          user.Email,
		"email":        user.Email,
		"fullName":     user.FullName,
		"phone_number": user.PhoneNumber,
		"iat":          time.Now().Unix(),
		"metadata": map[string]interface{}{
			"x-auth-allowed-roles": []string{user.Role},
			"x-auth-default-role":  user.Role,
			"x-auth-role":          user.Role,
			"x-auth-user-id":       user.Email,
		},
	}

	accessToken, err := jwt_auth.CreateAuthAccessToken(s.Config.Server.JWTKey, accessTokenPayload)
	if err != nil {
		return utils.SendErrorResponse(w, "failed to create access token", "token_creation_error", http.StatusInternalServerError)
	}
	responseData, _ := json.Marshal(payloads.AccessTokenOutput{
		AccessToken: accessToken,
	})

	if _, err := w.Write(responseData); err != nil {
		return err
	}

	return nil
}
