package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"symbiosisos/backend/internal/database" // Added the missing import!
)

// HandlerGetFacilities returns all physical locations owned by the user
func (apiCfg *APIConfig) HandlerGetFacilities(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value(UserIDKey).(string)
	parsedUserID, _ := uuid.Parse(userIDStr)

	facilities, err := apiCfg.DB.GetFacilitiesByUserID(r.Context(), uuid.NullUUID{
		UUID:  parsedUserID,
		Valid: true,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch facilities")
		return
	}

	if facilities == nil {
		facilities = []database.GetFacilitiesByUserIDRow{}
	}
	RespondWithJSON(w, http.StatusOK, facilities)
}

// HandlerGetFacilityWaste returns all active waste streams for a specific facility
func (apiCfg *APIConfig) HandlerGetFacilityWaste(w http.ResponseWriter, r *http.Request) {
	facilityIDStr := chi.URLParam(r, "facility_id")
	parsedFacilityID, err := uuid.Parse(facilityIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Facility ID")
		return
	}

	wasteStreams, err := apiCfg.DB.GetWasteStreamsByFacility(r.Context(), uuid.NullUUID{
		UUID:  parsedFacilityID,
		Valid: true,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch waste streams")
		return
	}

	if wasteStreams == nil {
		wasteStreams = []database.WasteStream{}
	}
	RespondWithJSON(w, http.StatusOK, wasteStreams)
}

// HandlerGetFacilityRequirements returns all active buy-orders for a specific facility
func (apiCfg *APIConfig) HandlerGetFacilityRequirements(w http.ResponseWriter, r *http.Request) {
	facilityIDStr := chi.URLParam(r, "facility_id")
	parsedFacilityID, err := uuid.Parse(facilityIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Facility ID")
		return
	}

	requirements, err := apiCfg.DB.GetRequirementsByFacility(r.Context(), uuid.NullUUID{
		UUID:  parsedFacilityID,
		Valid: true,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch requirements")
		return
	}

	if requirements == nil {
		requirements = []database.BuyerRequirement{}
	}
	RespondWithJSON(w, http.StatusOK, requirements)
}

// HandlerGetTransactions returns the user's ledger history
func (apiCfg *APIConfig) HandlerGetTransactions(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value(UserIDKey).(string)
	parsedUserID, _ := uuid.Parse(userIDStr)

	transactions, err := apiCfg.DB.GetTransactionsForUser(r.Context(), uuid.NullUUID{
		UUID:  parsedUserID,
		Valid: true,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch transaction history")
		return
	}

	if transactions == nil {
		transactions = []database.GetTransactionsForUserRow{}
	}
	RespondWithJSON(w, http.StatusOK, transactions)
}
