-- name: CreateUser :exec
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
    ?, ?, ?, ?
);

-- name: GetUser :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET hashed_password = ?
WHERE username = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = ?;