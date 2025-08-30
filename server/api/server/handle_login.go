package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	jwt_auth "github.com/c5rogers/one-tap/salary-advance-loan-system/internal/jwt-auth"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/internal/password"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) error {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	payload := models.LoginPayloadNew{}

	err = json.Unmarshal(reqBody, &payload)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	user, err := s.Graph.GetUserByEmail(payload.Email)
	if err != nil {
		return utils.SendErrorResponse(w, "error retriving user", "database_error", http.StatusInternalServerError)
	}

	if user == (&models.User{}) {
		return utils.SendErrorResponse(w, "username or password incorrect", "bad_credentials", http.StatusBadRequest)
	}

	if user == nil {
		return utils.SendErrorResponse(w, "username or password incorrect", "bad_credentials", http.StatusBadRequest)
	}

	passwordIsValid, msg := password.VerifyPassword(payload.Password, user.Password)
	if !passwordIsValid {
		return utils.SendErrorResponse(w, msg, "bad_credential", http.StatusBadRequest)
	}

	// fetch the roles here
	roles, err := s.Graph.GetAllowedRoles()
	if err != nil {
		return utils.SendErrorResponse(w, "error retriving allowed roles", "database_error", http.StatusInternalServerError)
	}

	if len(roles) == 0 {
		return utils.SendErrorResponse(w, "No allowed role found for user", "database_error", http.StatusInternalServerError)
	}

	allowedRoles := make([]string, len(roles))
	for i, role := range roles {
		allowedRoles[i] = role.Role
	}

	if len(user.Roles) == 0 {
		return utils.SendErrorResponse(w, "No assigned role for the user", "database_error", http.StatusInternalServerError)
	}

	tokenPayload := map[string]interface{}{
		"aud":          user.ID,
		"iss":          "https://hakaton.backend.com",
		"exp":          time.Now().Add(72 * time.Hour).Unix(),
		"sub":          user.ID,
		"email":        user.Email,
		"fullName":     user.FullName,
		"phone_number": user.PhoneNumber,
		"iat":          time.Now().Unix(),
		"metadata": map[string]interface{}{
			"x-hasura-allowed-roles": allowedRoles,
			"x-hasura-default-role":  user.Roles[0].Role,
			"x-hasura-role":          user.Roles[0].Role,
			"x-hasura-user-id":       user.ID,
		},
	}

	hasuraAccessToken, err := jwt_auth.CreateHasuraAccessToken(s.Config.Server.JWTKey, tokenPayload)
	if err != nil {
		return utils.SendErrorResponse(w, err.Error(), "token_error_sign_token", http.StatusInternalServerError)
	}

	responseData, _ := json.Marshal(models.AccessTokenOutput{
		AccessToken: hasuraAccessToken,
	})

	if _, err := w.Write(responseData); err != nil {
		return err
	}

	return nil
}
