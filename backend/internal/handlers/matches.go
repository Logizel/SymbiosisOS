package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// MatchResult is the final payload we send to the React frontend
type MatchResult struct {
	WasteStreamID      string  `json:"waste_stream_id"`
	BuyerRequirementID string  `json:"buyer_requirement_id"`
	BuyerName          string  `json:"buyer_name"`
	Chemical           string  `json:"chemical"`
	DistanceMiles      float64 `json:"distance_miles"`
	FreightCost        float64 `json:"freight_cost"`
	LandfillCost       float64 `json:"landfill_cost"`
	NetSavings         float64 `json:"net_savings"`
	IsViable           bool    `json:"is_viable"`
}

// HandlerGetMatches calculates the financial viability of spatial matches
func (apiCfg *APIConfig) HandlerGetMatches(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the Facility ID from the URL path (/api/v1/matches/{facility_id})
	facilityIDStr := chi.URLParam(r, "facility_id")
	facilityUUID, err := uuid.Parse(facilityIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Facility ID parameter")
		return
	}

	// 2. Run the PostGIS Database Query
	dbMatches, err := apiCfg.DB.GetViableMatchesForFacility(r.Context(), facilityUUID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to calculate spatial matches")
		return
	}

	// 3. The Arbitrage Gate (Financial Math)
	var finalMatches []MatchResult

	// Standard freight cost assumption: $5 per mile per truckload
	const FreightCostPerMile = 5.00 

	for _, match := range dbMatches {
		// Convert PostGIS meters to miles
		distanceMiles := match.DistanceMeters / 1609.34
		
		// Convert sqlc strings back to floats for math
		landfillFee, _ := strconv.ParseFloat(match.LocalLandfillFeePerTon, 64)
		
		// Calculate Costs
		freightCost := distanceMiles * FreightCostPerMile
		landfillCost := float64(match.TonnageAvailable) * landfillFee
		
		// Calculate Arbitrage (How much money they save by NOT going to the landfill)
		netSavings := landfillCost - freightCost
		
		// It is only viable if the savings are strictly positive
		isViable := netSavings > 0

		// We only want to return matches that actually save the factory money
		if isViable {
			finalMatches = append(finalMatches, MatchResult{
				WasteStreamID:      match.WasteStreamID.UUID.String(),
				BuyerRequirementID: match.BuyerRequirementID.UUID.String(),
				BuyerName:          match.BuyerName,
				Chemical:           match.PrimaryChemical,
				DistanceMiles:      math.Round(distanceMiles*100) / 100, // Round to 2 decimals
				FreightCost:        math.Round(freightCost*100) / 100,
				LandfillCost:       math.Round(landfillCost*100) / 100,
				NetSavings:         math.Round(netSavings*100) / 100,
				IsViable:           isViable,
			})
		}
	}

	// If the slice is empty, return an empty array instead of null
	if finalMatches == nil {
		finalMatches = []MatchResult{}
	}

	RespondWithJSON(w, http.StatusOK, finalMatches)
}
