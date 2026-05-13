package handlers

import (
	"database/sql" // Added for sql.NullString
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"symbiosisos/backend/internal/database"
)

// HandlerCreateWasteStream logs a specific chemical byproduct to the marketplace
func (apiCfg *APIConfig) HandlerCreateWasteStream(w http.ResponseWriter, r *http.Request) {
	// 1. Define the exact JSON structure from the frontend form
	type parameters struct {
		FacilityID             string  `json:"facility_id"`
		PrimaryChemical        string  `json:"primary_chemical"`
		PurityPercentage       float64 `json:"purity_percentage"`
		TonnageAvailable       int32   `json:"tonnage_available"`
		LocalLandfillFeePerTon float64 `json:"local_landfill_fee_per_ton"`
		LabReportURL           string  `json:"lab_report_url"`
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

	// 3. Format Decimals and Nullable Strings for sqlc
	// sqlc mapped your DECIMAL columns to standard Go strings to prevent float precision loss
	purityStr := fmt.Sprintf("%.2f", params.PurityPercentage)
	feeStr := fmt.Sprintf("%.2f", params.LocalLandfillFeePerTon)

	// sqlc mapped your nullable TEXT column to sql.NullString
	labURL := sql.NullString{
		String: params.LabReportURL,
		Valid:  params.LabReportURL != "", // Valid is true if the string is not empty
	}

	// 4. Execute the database insertion
	wasteStream, err := apiCfg.DB.CreateWasteStream(r.Context(), database.CreateWasteStreamParams{
		FacilityID: uuid.NullUUID{
			UUID:  parsedFacilityID,
			Valid: true,
		},
		PrimaryChemical:        params.PrimaryChemical,
		PurityPercentage:       purityStr, // Now passing it as the expected string
		TonnageAvailable:       params.TonnageAvailable,
		LocalLandfillFeePerTon: feeStr, // Now passing it as the expected string
		LabReportUrl:           labURL, // Now passing standard sql.NullString
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create waste stream")
		return
	}

	RespondWithJSON(w, http.StatusCreated, wasteStream)
}

func (apiCfg *APIConfig) HandlerDeleteWasteStream(w http.ResponseWriter, r *http.Request) {
	wasteIDStr := chi.URLParam(r, "id")
	wasteID, err := uuid.Parse(wasteIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid waste ID")
		return
	}

	// We also parse the facility ID from the query params to ensure ownership
	facilityIDStr := r.URL.Query().Get("facility_id")
	facilityID, err := uuid.Parse(facilityIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Facility ID required")
		return
	}

	err = apiCfg.DB.DeleteWasteStream(r.Context(), database.DeleteWasteStreamParams{
		ID:         wasteID,
		FacilityID: uuid.NullUUID{UUID: facilityID, Valid: true},
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete waste stream")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
