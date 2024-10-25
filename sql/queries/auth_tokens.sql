-- name: CreateAuthTokens :exec
INSERT INTO auth_tokens (user_id, access_token, refresh_token, expires_at)
VALUES (?, ?, ?, ?);


-- name: GetAuthTokens :one
SELECT user_id, access_token, refresh_token, expires_at, created_at
FROM auth_tokens
WHERE user_id = ?;

-- name: UpdateAuthTokens :exec
UPDATE auth_tokens
SET access_token = ?, refresh_token = ?, expires_at = ?
WHERE user_id = ?;

-- name: DeleteAuthTokens :exec
DELETE FROM auth_tokens
WHERE user_id = ?;