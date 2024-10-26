// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (membership_id, membership_type)
VALUES (?, ?)
RETURNING id, membership_id, membership_type, created_at
`

type CreateUserParams struct {
	MembershipID   string
	MembershipType int64
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.MembershipID, arg.MembershipType)
	var i User
	err := row.Scan(
		&i.ID,
		&i.MembershipID,
		&i.MembershipType,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, membership_id, membership_type, created_at
FROM users
WHERE id = ?
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.MembershipID,
		&i.MembershipType,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByMembershipID = `-- name: GetUserByMembershipID :one
SELECT id, membership_id, membership_type, created_at
FROM users
WHERE membership_id = ?
`

func (q *Queries) GetUserByMembershipID(ctx context.Context, membershipID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByMembershipID, membershipID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.MembershipID,
		&i.MembershipType,
		&i.CreatedAt,
	)
	return i, err
}
