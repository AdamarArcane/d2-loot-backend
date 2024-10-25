// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: auth_tokens.sql

package database

import (
	"context"
	"time"
)

const createAuthTokens = `-- name: CreateAuthTokens :exec
INSERT INTO auth_tokens (user_id, access_token, refresh_token, expires_at)
VALUES (?, ?, ?, ?)
`

type CreateAuthTokensParams struct {
	UserID       int64
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (q *Queries) CreateAuthTokens(ctx context.Context, arg CreateAuthTokensParams) error {
	_, err := q.db.ExecContext(ctx, createAuthTokens,
		arg.UserID,
		arg.AccessToken,
		arg.RefreshToken,
		arg.ExpiresAt,
	)
	return err
}

const deleteAuthTokens = `-- name: DeleteAuthTokens :exec
DELETE FROM auth_tokens
WHERE user_id = ?
`

func (q *Queries) DeleteAuthTokens(ctx context.Context, userID int64) error {
	_, err := q.db.ExecContext(ctx, deleteAuthTokens, userID)
	return err
}

const getAuthTokens = `-- name: GetAuthTokens :one
SELECT user_id, access_token, refresh_token, expires_at, created_at
FROM auth_tokens
WHERE user_id = ?
`

func (q *Queries) GetAuthTokens(ctx context.Context, userID int64) (AuthToken, error) {
	row := q.db.QueryRowContext(ctx, getAuthTokens, userID)
	var i AuthToken
	err := row.Scan(
		&i.UserID,
		&i.AccessToken,
		&i.RefreshToken,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateAuthTokens = `-- name: UpdateAuthTokens :exec
UPDATE auth_tokens
SET access_token = ?, refresh_token = ?, expires_at = ?
WHERE user_id = ?
`

type UpdateAuthTokensParams struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	UserID       int64
}

func (q *Queries) UpdateAuthTokens(ctx context.Context, arg UpdateAuthTokensParams) error {
	_, err := q.db.ExecContext(ctx, updateAuthTokens,
		arg.AccessToken,
		arg.RefreshToken,
		arg.ExpiresAt,
		arg.UserID,
	)
	return err
}
