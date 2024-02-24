-- name: CreateAccount :exec
INSERT INTO accounts (owner,
                      user_id,
                      balance,
                      currency, 
                      created_at)
VALUES ($1, $2, $3, 'USD', $4), ($1, $2, $3, 'EUR', $4), ($1, $2, $3, 'UAH', $4);

-- name: GetUserAccount :one
SELECT *
FROM accounts
WHERE user_id = $1 and id = $2;

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1;

-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE user_id = $1 and id = $2
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT *
FROM accounts
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddBalanceToAccount :one
UPDATE accounts SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;
