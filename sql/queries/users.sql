-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES (?, ?)
RETURNING id, email, created_at, updated_at;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: UpdateUserEmail :one
UPDATE users
SET
    email = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, email, created_at, updated_at;

-- name: UpdateUserPassword :one
UPDATE users
SET
    hashed_password = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, email, created_at, updated_at;
