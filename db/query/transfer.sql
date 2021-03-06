-- name: CreateTransfer :exec
INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
) VALUES (
    ?, ?, ?
);

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ?;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: GetTransferByFromAccount :many
SELECT * FROM transfers
WHERE from_account_id = ?
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: GetTransferByFromAccountAndToAccount :many
SELECT * FROM transfers
WHERE from_account_id = ? AND to_account_id = ?
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateTransfer :exec
UPDATE transfers
SET amount = ?
WHERE id = ?;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = ?;

-- name: GetLastTransfer :one
SELECT * FROM transfers
ORDER BY id DESC
LIMIT 1;