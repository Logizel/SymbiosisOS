-- name: GetFacilitiesByUserID :many
SELECT 
    id, 
    name, 
    ST_Y(coordinates::geometry)::float8 as lat, 
    ST_X(coordinates::geometry)::float8 as lng, 
    created_at
FROM facilities
WHERE user_id = $1;

-- name: GetWasteStreamsByFacility :many
SELECT * FROM waste_streams 
WHERE facility_id = $1
ORDER BY created_at DESC;

-- name: GetRequirementsByFacility :many
SELECT * FROM buyer_requirements 
WHERE facility_id = $1
ORDER BY created_at DESC;

-- name: GetTransactionsForUser :many
SELECT 
    t.id AS transaction_id,
    ws.primary_chemical,
    t.tonnage_exchanged,
    t.net_savings_estimated,
    t.status,
    t.created_at,
    f_gen.name AS generator_name,
    f_buy.name AS buyer_name
FROM transactions t
JOIN waste_streams ws ON t.waste_stream_id = ws.id
JOIN facilities f_gen ON ws.facility_id = f_gen.id
JOIN buyer_requirements br ON t.buyer_requirement_id = br.id
JOIN facilities f_buy ON br.facility_id = f_buy.id
WHERE f_gen.user_id = $1 OR f_buy.user_id = $1
ORDER BY t.created_at DESC;
