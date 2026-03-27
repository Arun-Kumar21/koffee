-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, device_info, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetRefresToken :one
SELECT id, user_id, token_hash, expires_at, revoked
FROM refresh_tokens
WHERE token_hash = $1 AND revoked = false;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = true
WHERE token_hash = $1;

-- name: RevokeAllUserTokens :exec
UPDATE refresh_tokens
SET revoked = true
WHERE user_id = $1;
