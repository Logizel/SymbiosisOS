package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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

func (apiCfg *APIConfig) HandlerGetMatches(w http.ResponseWriter, r *http.Request) {
	facilityIDStr := chi.URLParam(r, "facility_id")
	facilityUUID, err := uuid.Parse(facilityIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Facility ID parameter")
		return
	}

	dbMatches, err := apiCfg.DB.GetViableMatchesForFacility(r.Context(), facilityUUID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to calculate spatial matches")
		return
	}

	var finalMatches []MatchResult

	const FreightCostPerMile = 5.00

	for _, match := range dbMatches {
		distanceMiles := match.DistanceMeters / 1609.34

		landfillFee, _ := strconv.ParseFloat(match.LocalLandfillFeePerTon, 64)

		freightCost := distanceMiles * FreightCostPerMile
		landfillCost := float64(match.TonnageAvailable) * landfillFee

		netSavings := landfillCost - freightCost

		isViable := netSavings > 0

		if isViable {
			finalMatches = append(finalMatches, MatchResult{
				WasteStreamID:      match.WasteStreamID.String(),
				BuyerRequirementID: match.BuyerRequirementID.String(),
				BuyerName:          match.BuyerName,
				Chemical:           match.PrimaryChemical,
				DistanceMiles:      math.Round(distanceMiles*100) / 100,
				FreightCost:        math.Round(freightCost*100) / 100,
				LandfillCost:       math.Round(landfillCost*100) / 100,
				NetSavings:         math.Round(netSavings*100) / 100,
				IsViable:           isViable,
			})
		}
	}

	if finalMatches == nil {
		finalMatches = []MatchResult{}
	}

	RespondWithJSON(w, http.StatusOK, finalMatches)
}
