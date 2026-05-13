package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"symbiosisos/backend/internal/database"
)

func (apiCfg *APIConfig) HandlerCreateFacility(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value(UserIDKey).(string)

	parsedUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID format in token")
		return
	}

	type parameters struct {
		Name string  `json:"name"`
		Lat  float64 `json:"lat"`
		Lng  float64 `json:"lng"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	facility, err := apiCfg.DB.CreateFacility(r.Context(), database.CreateFacilityParams{
		UserID: uuid.NullUUID{
			UUID:  parsedUUID,
			Valid: true,
		},
		Name:    params.Name,
		Column3: params.Lng,
		Column4: params.Lat,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to map facility to the globe")
		return
	}

	RespondWithJSON(w, http.StatusCreated, facility)
}
