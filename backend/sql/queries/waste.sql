-- sql/queries/waste.sql

-- name: CreateWasteStream :one
INSERT INTO waste_streams (
    id, facility_id, primary_chemical, purity_percentage, 
    tonnage_available, local_landfill_fee_per_ton, lab_report_url
)
VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6
)
RETURNING *;
