-- sql/queries/facilities.sql

-- name: CreateFacility :one
INSERT INTO facilities (id, user_id, name, coordinates)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    ST_SetSRID(ST_MakePoint($3::float8, $4::float8), 4326)::geography
)
RETURNING id, user_id, name, 
          ST_Y(coordinates::geometry)::float8 as lat, 
          ST_X(coordinates::geometry)::float8 as lng, 
          created_at;
