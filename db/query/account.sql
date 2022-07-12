-- name: CreateAccount :exec
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
    ?, ?, ?
);

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = ?
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateAccount :exec
UPDATE accounts
SET balance = ?
WHERE id = ?;

-- name: AddAccountBalance :exec
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id);

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;

-- name: GetLastAccount :one
SELECT * FROM accounts
ORDER BY id DESC
LIMIT 1;