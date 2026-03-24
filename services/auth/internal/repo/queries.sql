-- name: CreateUser :one
INSERT INTO users (name, email, password, avatar_url)
VALUES ($1, $2, $3, $4)
RETURNING id, email, created_at;


-- name: GetUserByEmail :one
SELECT name, email, avatar_url from users
WHERE email = $1;

