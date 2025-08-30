package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	jwt_auth "github.com/c5rogers/one-tap/salary-advance-loan-system/internal/jwt-auth"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
	"github.com/golang-jwt/jwt/v5"
)

// func (s *Server) VerifyActionSecret(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		// Extract and verify the ACTION_SECRET_KEY header
// 		actionSecret := r.Header.Get("ACTION_SECRET")
// 		if actionSecret != s.Config.Server.ActionSecretKey {
// 			utils.SendErrorResponse(w, "Invalid action secret key", "invalid_action_secret_key", http.StatusUnauthorized)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func (s *Server) ExtractUserFromToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		publicKey, err := jwt_auth.LoadECDSAPublicKey(s.Config.Server.JWTPublicKey)
		if err != nil {
			utils.SendErrorResponse(w, "Failed to load public key", "invalid_public_key", http.StatusInternalServerError)
			return
		}

		claims := &models.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			utils.SendErrorResponse(w, "Invalid token", "invalid_token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
