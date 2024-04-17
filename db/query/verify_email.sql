-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (user_id,
                      email,
                      code,
                      is_used,
                      expired_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (email) DO UPDATE SET code = $3
RETURNING *;

-- name: GetVerifyEmail :one
SELECT *
FROM verify_emails
WHERE id = $1;

-- name: UpdateVerifyEmails :exec
UPDATE verify_emails
SET is_used = $1
WHERE id = $2;
