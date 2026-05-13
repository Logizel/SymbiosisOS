package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"symbiosisos/backend/internal/database"
)

func (apiCfg *APIConfig) HandlerCreateRequirement(w http.ResponseWriter, r *http.Request) {
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

	parsedFacilityID, err := uuid.Parse(params.FacilityID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Facility ID")
		return
	}

	purityStr := fmt.Sprintf("%.2f", params.MinimumPurity)

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

func (apiCfg *APIConfig) HandlerDeleteRequirement(w http.ResponseWriter, r *http.Request) {
	reqIDStr := chi.URLParam(r, "id")
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid requirement ID")
		return
	}

	facilityIDStr := r.URL.Query().Get("facility_id")
	facilityID, err := uuid.Parse(facilityIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Facility ID required")
		return
	}

	err = apiCfg.DB.DeleteBuyerRequirement(r.Context(), database.DeleteBuyerRequirementParams{
		ID:         reqID,
		FacilityID: uuid.NullUUID{UUID: facilityID, Valid: true},
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete requirement")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
