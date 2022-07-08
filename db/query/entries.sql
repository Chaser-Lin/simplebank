-- name: CreateEntry :exec
INSERT INTO entries (
    account_id,
    amount
)
VALUES(
     ?, ?
);

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ?;

-- name: GetEntryByAccount :many
SELECT * FROM entries
WHERE account_id = ?
LIMIT ?
OFFSET ?;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateEntry :exec
UPDATE entries
SET amount = ?
WHERE id = ?;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = ?;

-- name: GetLastEntry :one
SELECT * FROM entries
ORDER BY id DESC
LIMIT 1;