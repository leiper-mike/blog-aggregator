-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;
-- name: GetUser :one
SELECT * 
FROM users
WHERE $1 = name 
LIMIT 1;
-- name: GetUserById :one
SELECT * 
FROM users
WHERE $1 = id
LIMIT 1;
-- name: DeleteAll :exec
DELETE FROM users;
-- name: GetAll :many
SELECT * FROM users;