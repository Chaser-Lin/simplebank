package db

import (
	"SimpleBank/db/util"
	"context"
	_ "database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

//func  (q *Queries) CreateAndReturnAccount(ctx context.Context, arg CreateAccountParams) (result Account, err error) {
//	err = q.CreateAccount(ctx, arg)
//	if err != nil {
//		return
//	}
//	result, err = q.GetLastAccount(ctx)
//	return
//}

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAndReturnAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.ID, account1.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account1.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestQueries_ListAccounts(t *testing.T) {
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 3,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	for _, account1 := range accounts {
		//fmt.Println(account1)
		account2, err := testQueries.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)
		require.Equal(t, account1.ID, account2.ID)
		require.Equal(t, account1.Owner, account2.Owner)
		require.Equal(t, account1.Balance, account2.Balance)
		require.Equal(t, account1.Currency, account2.Currency)
		require.Equal(t, account1.CreatedAt, account2.CreatedAt)
	}
}

func TestQueries_UpdateAccount(t *testing.T) {
	arg := UpdateAccountParams{
		Balance: 999.99,
		ID:      2,
	}
	testQueries.UpdateAccount(context.Background(), arg)
	account, err := testQueries.GetAccount(context.Background(), arg.ID)
	require.NoError(t, err)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.ID, arg.ID)
}

func TestQueries_DeleteAccount(t *testing.T) {
	var id int64 = 11
	testQueries.DeleteAccount(context.Background(), id)
	accout, err := testQueries.GetAccount(context.Background(), id)
	require.Error(t, err)
	require.Empty(t, accout)
}
