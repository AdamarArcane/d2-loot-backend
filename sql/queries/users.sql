-- name: GetUserByMembershipID :one
SELECT id, membership_id, membership_type, created_at
FROM users
WHERE membership_id = ?;

-- name: CreateUser :one
INSERT INTO users (membership_id, membership_type)
VALUES (?, ?)
RETURNING id, membership_id, membership_type, created_at;

-- name: GetUser :one
SELECT id, membership_id, membership_type, created_at
FROM users
WHERE id = ?;
