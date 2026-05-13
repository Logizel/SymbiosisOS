-- sql/queries/delete.sql

-- name: DeleteWasteStream :exec
DELETE FROM waste_streams
WHERE id = $1 AND facility_id = $2;

-- name: DeleteBuyerRequirement :exec
DELETE FROM buyer_requirements
WHERE id = $1 AND facility_id = $2;
