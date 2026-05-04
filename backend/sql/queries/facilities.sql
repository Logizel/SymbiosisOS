-- sql/queries/facilities.sql

-- name: CreateFacility :one
INSERT INTO facilities (id, user_id, name, coordinates)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    -- We take Longitude ($3) and Latitude ($4), turn them into a geometric point, 
    -- set the spatial reference to 4326 (WGS 84 / GPS Standard), and cast to geography.
    ST_SetSRID(ST_MakePoint($3::float8, $4::float8), 4326)::geography
)
-- Instead of returning the raw binary geography data (which Go will struggle to read), 
-- we use PostGIS functions to extract the exact Lat/Lng floats back out.
RETURNING id, user_id, name, 
          ST_Y(coordinates::geometry)::float8 as lat, 
          ST_X(coordinates::geometry)::float8 as lng, 
          created_at;
