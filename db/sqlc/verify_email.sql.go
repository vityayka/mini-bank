// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: verify_email.sql

package db

import (
	"context"
	"time"
)

const createVerifyEmail = `-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (user_id,
                      email,
                      code,
                      is_used,
                      expired_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (email) DO UPDATE SET code = $3
RETURNING id, user_id, email, code, is_used, created_at, expired_at
`

type CreateVerifyEmailParams struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	IsUsed    bool      `json:"is_used"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, createVerifyEmail,
		arg.UserID,
		arg.Email,
		arg.Code,
		arg.IsUsed,
		arg.ExpiredAt,
	)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Email,
		&i.Code,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const getVerifyEmail = `-- name: GetVerifyEmail :one
SELECT id, user_id, email, code, is_used, created_at, expired_at
FROM verify_emails
WHERE id = $1
`

func (q *Queries) GetVerifyEmail(ctx context.Context, id int64) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, getVerifyEmail, id)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Email,
		&i.Code,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const updateVerifyEmails = `-- name: UpdateVerifyEmails :exec
UPDATE verify_emails
SET is_used = $1
WHERE id = $2
`

type UpdateVerifyEmailsParams struct {
	IsUsed bool  `json:"is_used"`
	ID     int64 `json:"id"`
}

func (q *Queries) UpdateVerifyEmails(ctx context.Context, arg UpdateVerifyEmailsParams) error {
	_, err := q.db.Exec(ctx, updateVerifyEmails, arg.IsUsed, arg.ID)
	return err
}
