package handlers

import (
	"encoding/json"
	"net/http"

	"symbiosisos/backend/internal/database"
)

type APIConfig struct {
	DB *database.Queries
}

// HandlerCreateUser processes a POST request to create a new user
func (apiCfg *APIConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Define the exact JSON structure we expect from the React frontend
	type parameters struct {
		Email       string `json:"email"`
		Password    string `json:"password"` // Note: In a production app, we will hash this!
		Role        string `json:"role"`
		CompanyName string `json:"company_name"`
	}

	// 2. Decode the incoming JSON body into our struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON format") // Capital R
		return
	}

	// 3. Call the sqlc generated database method
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:        params.Email,
		PasswordHash: params.Password,
		Role:         params.Role,
		CompanyName:  params.CompanyName,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create user in database") // Capital R
		return
	}

	// 4. Return the newly created user as JSON
	RespondWithJSON(w, http.StatusCreated, user) // Capital R
}
