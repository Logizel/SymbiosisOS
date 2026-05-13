package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"symbiosisos/backend/internal/database"
)

func (apiCfg *APIConfig) HandlerCreateWasteStream(w http.ResponseWriter, r *http.Request) {
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

	parsedFacilityID, err := uuid.Parse(params.FacilityID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Facility ID")
		return
	}

	purityStr := fmt.Sprintf("%.2f", params.PurityPercentage)
	feeStr := fmt.Sprintf("%.2f", params.LocalLandfillFeePerTon)

	labURL := sql.NullString{
		String: params.LabReportURL,
		Valid:  params.LabReportURL != "",
	}

	wasteStream, err := apiCfg.DB.CreateWasteStream(r.Context(), database.CreateWasteStreamParams{
		FacilityID: uuid.NullUUID{
			UUID:  parsedFacilityID,
			Valid: true,
		},
		PrimaryChemical:        params.PrimaryChemical,
		PurityPercentage:       purityStr,
		TonnageAvailable:       params.TonnageAvailable,
		LocalLandfillFeePerTon: feeStr,
		LabReportUrl:           labURL,
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
