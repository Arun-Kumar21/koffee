-- name: CreateUser :one
INSERT INTO users (name, email, password, avatar_url)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, avatar_url, role, created_at;

-- name: GetUserByEmail :one
SELECT id, name, email, password, avatar_url, role, created_at from users 
WHERE email = $1;

-- name: UpdateUserProfile :one 
UPDATE users
SET name = $1,
    avatar_url = $2,
    updated_at = NOW()
where email = $3
RETURNING id, name, email, avatar_url, role, created_at;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1,
    updated_at = NOW()
WHERE id = $2;