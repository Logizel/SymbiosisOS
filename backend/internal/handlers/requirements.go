package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"symbiosisos/backend/internal/database"
	"github.com/google/uuid"
)

// HandlerCreateRequirement logs a buyer's standing order to the marketplace
func (apiCfg *APIConfig) HandlerCreateRequirement(w http.ResponseWriter, r *http.Request) {
	// 1. Define the exact JSON structure we expect from the frontend
	type parameters struct {
		FacilityID                  string  `json:"facility_id"`
		RequiredChemical            string  `json:"required_chemical"`
		MinimumPurity               float64 `json:"minimum_purity"`
		MaxAcceptableDistanceMeters int32   `json:"max_acceptable_distance_meters"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// 2. Parse the Facility UUID
	parsedFacilityID, err := uuid.Parse(params.FacilityID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Facility ID")
		return
	}

	// 3. Format Decimal for sqlc (Just like we did for the waste stream)
	purityStr := fmt.Sprintf("%.2f", params.MinimumPurity)

	// 4. Execute the database insertion
	requirement, err := apiCfg.DB.CreateBuyerRequirement(r.Context(), database.CreateBuyerRequirementParams{
		FacilityID: uuid.NullUUID{
			UUID:  parsedFacilityID,
			Valid: true,
		},
		RequiredChemical:            params.RequiredChemical,
		MinimumPurity:               purityStr, 
		MaxAcceptableDistanceMeters: params.MaxAcceptableDistanceMeters,
	})
	
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create buyer requirement")
		return
	}

	RespondWithJSON(w, http.StatusCreated, requirement)
}
