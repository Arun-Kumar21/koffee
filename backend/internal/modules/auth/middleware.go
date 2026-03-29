package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type contextKey string

const UserIDkey  contextKey = "userID"


func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len (parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return 
		}

		claims, err := h.tokens.ValidateAccessToken(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDkey, claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func GetUserID (r* http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDkey).(string)
	if !ok {
		return "", errors.New("user not authenticated")
	}

	return userID, nil
}
