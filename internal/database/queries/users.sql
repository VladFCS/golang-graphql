-- name: ListUsers :many
SELECT
    id,
    email,
    username,
    password_hash,
    created_at,
    updated_at
FROM users
ORDER BY created_at DESC, id DESC;

-- name: GetUserByID :one
SELECT
    id,
    email,
    username,
    password_hash,
    created_at,
    updated_at
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    password_hash
) VALUES (
    $1,
    $2,
    $3
)
RETURNING
    id,
    email,
    username,
    password_hash,
    created_at,
    updated_at;
