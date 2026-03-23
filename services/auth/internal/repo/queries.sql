-- name: CreateUser :one
INSERT INTO users (first_name, email, password, last_name)
VALUES ($1, $2, $3, $4)
RETURNING id, email, created_at;


