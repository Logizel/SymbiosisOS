package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"symbiosisos/backend/internal/database"
)

// HandlerCreateFacility registers a physical location for the authenticated user
func (apiCfg *APIConfig) HandlerCreateFacility(w http.ResponseWriter, r *http.Request) {
	// 1. Get the authenticated User ID from our JWT Middleware context
	userIDStr := r.Context().Value(UserIDKey).(string)

	// Convert the string ID into a Google UUID
	parsedUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID format in token")
		return
	}

	// 2. Define the JSON structure we expect from the frontend
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

	// 3. Execute the PostGIS database query
	// Note: We wrap the parsedUUID in a NullUUID struct and mark it Valid: true
	facility, err := apiCfg.DB.CreateFacility(r.Context(), database.CreateFacilityParams{
		UserID: uuid.NullUUID{
			UUID:  parsedUUID,
			Valid: true,
		},
		Name:    params.Name,
		Column3: params.Lng, // $3 in our SQL was Longitude
		Column4: params.Lat, // $4 in our SQL was Latitude
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to map facility to the globe")
		return
	}

	// 4. Return the successfully mapped facility
	RespondWithJSON(w, http.StatusCreated, facility)
}
