package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"symbiosisos/backend/internal/auth"
	"symbiosisos/backend/internal/database"
)

// APIConfig holds our database connection state
type APIConfig struct {
	DB *database.Queries
}

// HandlerCreateUser processes a POST request to create a new user
func (apiCfg *APIConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		Role        string `json:"role"`
		CompanyName string `json:"company_name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Call the sqlc generated database method
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:        params.Email,
		PasswordHash: params.Password,
		Role:         params.Role,
		CompanyName:  params.CompanyName,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create user in database")
		return
	}

	RespondWithJSON(w, http.StatusCreated, user)
}

// HandlerLogin verifies credentials and returns a JWT
func (apiCfg *APIConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// 1. Find the user in the database by email
	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// 2. Check the password
	if user.PasswordHash != params.Password {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// 3. Generate the JWT (Valid for 24 hours)
	tokenSecret := "development_super_secret_key"
	userIDStr := fmt.Sprintf("%v", user.ID)

	token, err := auth.MakeJWT(userIDStr, tokenSecret, time.Hour*24)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	// 4. Return the token to the user
	RespondWithJSON(w, http.StatusOK, map[string]string{
		"id":    userIDStr,
		"email": user.Email,
		"token": token,
	})
}

// HandlerGetMe returns the authenticated user's ID to prove the middleware works
func (apiCfg *APIConfig) HandlerGetMe(w http.ResponseWriter, r *http.Request) {
	// Extract the User ID that our middleware injected into the context
	userID := r.Context().Value(UserIDKey).(string)

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "You are securely authenticated!",
		"user_id": userID,
	})
}
