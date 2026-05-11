package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"symbiosisos/backend/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *APIConfig) HandlerCreateTransaction(w http.ResponseWriter, r *http.Request) {
	// 1. We expect the exact data outputted from our Matches engine
	type parameters struct {
		WasteStreamID      string  `json:"waste_stream_id"`
		BuyerRequirementID string  `json:"buyer_requirement_id"`
		TonnageExchanged   int32   `json:"tonnage_exchanged"`
		FreightCost        float64 `json:"freight_cost"`
		NetSavings         float64 `json:"net_savings"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// 2. Parse the IDs
	wasteID, err := uuid.Parse(params.WasteStreamID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Waste Stream ID")
		return
	}

	buyerID, err := uuid.Parse(params.BuyerRequirementID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Buyer Requirement ID")
		return
	}

	// 3. Format Decimals for sqlc
	freightStr := fmt.Sprintf("%.2f", params.FreightCost)
	savingsStr := fmt.Sprintf("%.2f", params.NetSavings)

	// 4. Execute the Atomic CTE Query
	transaction, err := apiCfg.DB.ExecuteMatchTransaction(r.Context(), database.ExecuteMatchTransactionParams{
		TonnageAvailable:      params.TonnageExchanged, // The $1 parameter in our SQL
		ID:                    wasteID,                 // The $2 parameter
		BuyerRequirementID:    buyerID,                 // The $3 parameter
		FreightCostEstimated:  freightStr,              // The $4 parameter
		NetSavingsEstimated:   savingsStr,              // The $5 parameter
	})
	
	if err != nil {
		// If the CTE fails (usually because of insufficient tonnage), sqlc throws a NoRows error
		RespondWithError(w, http.StatusConflict, "Transaction failed: Insufficient tonnage available or stream no longer exists.")
		return
	}

	// 5. Return the finalized digital contract
	RespondWithJSON(w, http.StatusCreated, transaction)
}
