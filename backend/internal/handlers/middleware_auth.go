package handlers

import (
	"context"
	"net/http"

	"symbiosisos/backend/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

func (apiCfg *APIConfig) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Missing or invalid auth header")
			return
		}

		tokenSecret := "development_super_secret_key"
		userID, err := auth.ValidateJWT(tokenString, tokenSecret)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
