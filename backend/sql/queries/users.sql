-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, role, company_name)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
