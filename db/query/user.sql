-- name: CreateUser :one
INSERT INTO users (username,
                      hashed_password,
                      full_name,
                      email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email),
    is_verified = COALESCE(sqlc.narg(is_verified), is_verified),
    password_changed_at = (CASE WHEN sqlc.narg(hashed_password) IS null THEN password_changed_at ELSE now() END)
WHERE id = sqlc.arg(id)
RETURNING *;