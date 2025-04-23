-- name: CreateUser :one
INSERT INTO users (name, username, email, hashed_password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1;
