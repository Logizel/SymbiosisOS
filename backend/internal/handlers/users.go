package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"symbiosisos/backend/internal/auth"
	"symbiosisos/backend/internal/database"
)

type APIConfig struct {
	DB *database.Queries
}

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

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if user.PasswordHash != params.Password {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	tokenSecret := "development_super_secret_key"
	userIDStr := fmt.Sprintf("%v", user.ID)

	token, err := auth.MakeJWT(userIDStr, tokenSecret, time.Hour*24)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"id":    userIDStr,
		"email": user.Email,
		"token": token,
	})
}

func (apiCfg *APIConfig) HandlerGetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "You are securely authenticated!",
		"user_id": userID,
	})
}
