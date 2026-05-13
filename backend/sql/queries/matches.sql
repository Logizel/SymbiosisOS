-- sql/queries/matches.sql

-- name: GetViableMatchesForFacility :many
SELECT 
    ws.id AS waste_stream_id,
    ws.primary_chemical,
    ws.tonnage_available,
    ws.local_landfill_fee_per_ton,
    br.id AS buyer_requirement_id,
    f_buyer.name AS buyer_name,
    -- Calculate the exact distance in meters over the Earth's curvature
    ST_Distance(f_generator.coordinates, f_buyer.coordinates)::float8 AS distance_meters
FROM waste_streams ws
JOIN facilities f_generator ON ws.facility_id = f_generator.id
-- The Chemical Gate: Exact string match
JOIN buyer_requirements br ON ws.primary_chemical = br.required_chemical
JOIN facilities f_buyer ON br.facility_id = f_buyer.id
WHERE f_generator.id = $1
  AND ws.tonnage_available > 0
  AND ws.purity_percentage >= br.minimum_purity
  AND ST_DWithin(f_generator.coordinates, f_buyer.coordinates, br.max_acceptable_distance_meters);
