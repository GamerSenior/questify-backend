package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/GamerSenior/questify-backend/models"
	u "github.com/GamerSenior/questify-backend/utils"
	"github.com/dgrijalva/jwt-go"
)

func writeErrorResponse(w http.ResponseWriter, message string) {
	response := u.Message(false, message)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, response)
}

// JwtAuthentication explanation
var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/users/new", "/api/users/login"}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			writeErrorResponse(w, "Missing auth token")
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			writeErrorResponse(w, "Invalid/Malformed auth token")
			return
		}

		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			writeErrorResponse(w, "Malformed authentication token")
			return
		}

		if !token.Valid {
			writeErrorResponse(w, "Token is not valid")
			return
		}

		fmt.Printf("User: %v", tk.UserID)
		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
