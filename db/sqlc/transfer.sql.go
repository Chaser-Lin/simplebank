// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: transfer.sql

package db

import (
	"context"
)

const createTransfer = `-- name: CreateTransfer :exec
INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
) VALUES (
    ?, ?, ?
)
`

type CreateTransferParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) error {
	_, err := q.db.ExecContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	return err
}

const deleteTransfer = `-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = ?
`

func (q *Queries) DeleteTransfer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransfer, id)
	return err
}

const getLastTransfer = `-- name: GetLastTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
ORDER BY id DESC
LIMIT 1
`

func (q *Queries) GetLastTransfer(ctx context.Context) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getLastTransfer)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE id = ?
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransferByFromAccount = `-- name: GetTransferByFromAccount :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE from_account_id = ?
ORDER BY id
LIMIT ?
OFFSET ?
`

type GetTransferByFromAccountParams struct {
	FromAccountID int64 `json:"from_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (q *Queries) GetTransferByFromAccount(ctx context.Context, arg GetTransferByFromAccountParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, getTransferByFromAccount, arg.FromAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransferByFromAccountAndToAccount = `-- name: GetTransferByFromAccountAndToAccount :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE from_account_id = ? AND to_account_id = ?
ORDER BY id
LIMIT ?
OFFSET ?
`

type GetTransferByFromAccountAndToAccountParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (q *Queries) GetTransferByFromAccountAndToAccount(ctx context.Context, arg GetTransferByFromAccountAndToAccountParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, getTransferByFromAccountAndToAccount,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
ORDER BY id
LIMIT ?
OFFSET ?
`

type ListTransfersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTransfer = `-- name: UpdateTransfer :exec
UPDATE transfers
SET amount = ?
WHERE id = ?
`

type UpdateTransferParams struct {
	Amount float64 `json:"amount"`
	ID     int64   `json:"id"`
}

func (q *Queries) UpdateTransfer(ctx context.Context, arg UpdateTransferParams) error {
	_, err := q.db.ExecContext(ctx, updateTransfer, arg.Amount, arg.ID)
	return err
}

func  (q *Queries) CreateAndReturnTransfer(ctx context.Context, arg CreateTransferParams) (result Transfer, err error) {
	err = q.CreateTransfer(ctx, arg)
	if err != nil {
		return
	}
	result, err = q.GetLastTransfer(ctx)
	return
}