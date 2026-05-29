-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES ($1, $2)
RETURNING id, email, created_at, updated_at;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserEmail :one
UPDATE users
SET
    email = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, created_at, updated_at;

-- name: UpdateUserPassword :one
UPDATE users
SET
    hashed_password = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, created_at, updated_at;
