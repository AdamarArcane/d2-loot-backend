// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"
)

type AuthToken struct {
	UserID       int64
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    sql.NullTime
}

type User struct {
	ID             int64
	MembershipID   string
	MembershipType int64
	CreatedAt      sql.NullTime
}
