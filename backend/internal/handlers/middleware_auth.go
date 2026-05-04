package handlers

import (
	"context"
	"net/http"

	"symbiosisos/backend/internal/auth"
)

// Define a custom type for our context keys to prevent collisions
type contextKey string
const UserIDKey contextKey = "userID"

// MiddlewareAuth intercepts requests to ensure a valid JWT is present
func (apiCfg *APIConfig) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract the token from the headers
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Missing or invalid auth header")
			return
		}

		// 2. Validate the token (Must match the secret used in HandlerLogin!)
		tokenSecret := "development_super_secret_key"
		userID, err := auth.ValidateJWT(tokenString, tokenSecret)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// 3. Store the User ID in the request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// 4. Pass the modified request to the actual route handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
