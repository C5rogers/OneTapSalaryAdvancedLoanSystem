package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	jwt_auth "github.com/c5rogers/one-tap/salary-advance-loan-system/internal/jwt-auth"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/internal/password"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
)

var validate = validator.New()

func (s *Server) HandleSignup(w http.ResponseWriter, r *http.Request) error {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	payload := models.SignUpPayload{}
	err = json.Unmarshal(reqBody, &payload)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	//validate the payload
	validationError := validate.Struct(payload)
	if validationError != nil {
		return utils.SendErrorResponse(w, "invalid SignUp information", "invalid_payload", http.StatusBadRequest)
	}

	user, err := s.Graph.GetUserByEmail(payload.Email)
	if err != nil {
		return utils.SendErrorResponse(w, "error retriving user", "database_error", http.StatusInternalServerError)
	}

	if user != nil {
		return utils.SendErrorResponse(w, "user already exists", "user_already_exists", http.StatusBadRequest)
	}
	// has the password
	hashedPassword, err := password.HashPassword(payload.Password)

	registeringUser := models.RegisteringUser{
		FullName:       payload.FullName,
		Email:          payload.Email,
		PhoneNumber:    payload.PhoneNumber,
		Password:       hashedPassword,
		ProfilePicture: payload.ProfilePicture,
		Role:           "user",
	}

	registeredUser, err := s.Graph.RegisterUser(&registeringUser)
	if err != nil {
		return utils.SendErrorResponse(w, "error registering user", "database_error", http.StatusInternalServerError)
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

	tokenPayload := map[string]interface{}{
		"aud":          registeredUser.ID,
		"iss":          "https://hakaton.backend.com",
		"exp":          time.Now().Add(72 * time.Hour).Unix(),
		"sub":          registeredUser.ID,
		"email":        registeredUser.Email,
		"fullName":     registeredUser.FullName,
		"phone_number": registeredUser.PhoneNumber,
		"iat":          time.Now().Unix(),
		"metadata": map[string]interface{}{
			"x-hasura-allowed-roles": allowedRoles,
			"x-hasura-default-role":  registeredUser.Roles[0].Role,
			"x-hasura-role":          registeredUser.Roles[0].Role,
			"x-hasura-user-id":       registeredUser.ID,
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
