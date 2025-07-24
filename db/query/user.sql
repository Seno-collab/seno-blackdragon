-- name CreateUser :one
INSERT INTO users (name) VALUES ($1) RETURNING id, name;

-- name GetUserByID :one
SELECT id, name FROM users WHERE id = $1;
