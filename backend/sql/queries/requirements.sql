-- sql/queries/requirements.sql

-- name: CreateBuyerRequirement :one
INSERT INTO buyer_requirements (
    id, facility_id, required_chemical, minimum_purity, max_acceptable_distance_meters
)
VALUES (
    gen_random_uuid(), $1, $2, $3, $4
)
RETURNING *;
