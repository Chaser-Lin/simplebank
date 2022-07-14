package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrInsufficientBalance = errors.New("account's balance is insufficient")

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParms) (TransferTxResult, error)
	WithdrawTx(ctx context.Context, arg BusinessTxParms) (BusinessTxResult, error)
	DepositTx(ctx context.Context, arg BusinessTxParms) (BusinessTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	//开始事务
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	//创建该事务用的Queries
	q := New(tx)
	//执行事务
	err = fn(q)
	if err != nil {
		//执行出错就rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	//提交事务，把错误返回
	return tx.Commit()
}

//包含Transfer事务的参数
type TransferTxParms struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParms) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) (err error) {
		result.Transfer, err = q.CreateAndReturnTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return
		}

		result.FromEntry, err = q.CreateAndReturnEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return
		}
		result.ToEntry, err = q.CreateAndReturnEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return
		}

		if arg.FromAccountID < arg.ToAccountID {
			//result.FromAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount)
			//if err != nil {
			//	return
			//}
			//result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount)
			result.FromAccount, result.ToAccount, err = transferMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			//result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount)
			//if err != nil {
			//	return
			//}
			//result.FromAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount)
			result.ToAccount, result.FromAccount, err = transferMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return
	})
	return result, err
}

type BusinessTxParms struct {
	AccountID int64   `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type BusinessTxResult struct {
	Entry   Entry   `json:"entry"`
	Account Account `json:"account"`
}

func (s *SQLStore) WithdrawTx(ctx context.Context, arg BusinessTxParms) (BusinessTxResult, error) {
	var result BusinessTxResult
	err := s.execTx(ctx, func(q *Queries) (err error) {
		result.Entry, err = q.CreateAndReturnEntry(ctx, CreateEntryParams{
			AccountID: arg.AccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return
		}
		result.Account, err = addMoney(ctx, q, arg.AccountID, -arg.Amount)
		return
	})
	return result, err
}

func (s *SQLStore) DepositTx(ctx context.Context, arg BusinessTxParms) (BusinessTxResult, error) {
	var result BusinessTxResult
	err := s.execTx(ctx, func(q *Queries) (err error) {
		result.Entry, err = q.CreateAndReturnEntry(ctx, CreateEntryParams{
			AccountID: arg.AccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return
		}
		result.Account, err = addMoney(ctx, q, arg.AccountID, arg.Amount)
		return
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID int64,
	amount float64,
) (account Account, err error) {
	account, err = q.GetAccount(ctx, accountID)
	if err != nil {
		return
	}
	err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount,
		ID:     accountID,
	})
	if err != nil {
		return
	}
	account, err = q.GetAccount(ctx, accountID)
	if err != nil {
		return
	}
	//操作之后余额小于0，出错回滚
	if account.Balance < 0 {
		err = ErrInsufficientBalance
	}
	return
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 float64,
	accountID2 int64,
	amount2 float64,
) (account1 Account, account2 Account, err error) {
	err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}
	account1, err = q.GetAccount(ctx, accountID1)
	if err != nil {
		return
	}
	err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}
	account2, err = q.GetAccount(ctx, accountID2)
	if err != nil {
		return
	}
	return
}
