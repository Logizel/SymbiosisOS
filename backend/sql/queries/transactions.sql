-- sql/queries/transactions.sql

-- name: ExecuteMatchTransaction :one
WITH updated_waste AS (
    -- Step A: Attempt to deduct the tonnage
    UPDATE waste_streams
    SET tonnage_available = tonnage_available - $1
    WHERE id = $2 AND tonnage_available >= $1
    RETURNING id
)
-- Step B: Only if Step A succeeds, create the ledger entry
INSERT INTO transactions (
    id, waste_stream_id, buyer_requirement_id, tonnage_exchanged, 
    freight_cost_estimated, net_savings_estimated
)
SELECT 
    gen_random_uuid(), $2, $3, $1, $4, $5
WHERE EXISTS (SELECT 1 FROM updated_waste)
RETURNING *;
