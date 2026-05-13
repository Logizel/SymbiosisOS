package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"symbiosisos/backend/internal/database"
)

func (apiCfg *APIConfig) HandlerCreateTransaction(w http.ResponseWriter, r *http.Request) {
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

	freightStr := fmt.Sprintf("%.2f", params.FreightCost)
	savingsStr := fmt.Sprintf("%.2f", params.NetSavings)

	transaction, err := apiCfg.DB.ExecuteMatchTransaction(r.Context(), database.ExecuteMatchTransactionParams{
		TonnageExchanged:     params.TonnageExchanged,
		WasteStreamID:        wasteID,
		BuyerRequirementID:   buyerID,
		FreightCostEstimated: freightStr,
		NetSavingsEstimated:  savingsStr,
	})

	if err != nil {
		RespondWithError(w, http.StatusConflict, "Transaction failed: Insufficient tonnage available or stream no longer exists.")
		return
	}

	RespondWithJSON(w, http.StatusCreated, transaction)
}
