package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	n := 5
	amount := float64(100.00)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParms{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		//check Transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check Entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//get account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		//check accounts' balance diffrence
		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := int(account1.Balance - fromAccount.Balance)
		diff2 := int(toAccount.Balance - account2.Balance)
		//fmt.Println(diff1, " ", diff2)
		require.True(t, diff1 == diff2)
		require.True(t, diff1 > 0)

		k := diff1 / int(amount)
		require.True(t, k >= 1 && k <= n)
		//fmt.Println("k = ", k)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-float64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+float64(n)*amount, updateAccount2.Balance)
}

func TestStore_TransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	n := 10
	amount := float64(10.00)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID, toAccountID = toAccountID, fromAccountID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParms{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}

func TestSQLStore_WithdrawTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance)

	n := 5
	amount := float64(10.00)

	errs := make(chan error, 5)
	results := make(chan BusinessTxResult, 5)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.WithdrawTx(context.Background(), BusinessTxParms{
				AccountID: account1.ID,
				Amount:    amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	m := 0
	for i := 0; i < n; i++ {
		result := <-results
		require.NotEmpty(t, result)
		err := <-errs
		if err != nil {
			require.True(t, result.Account.Balance-amount < 0)
			break
		}
		require.NoError(t, err)
		m++

		//check Entries
		entry := result.Entry
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account1.ID)
		require.Equal(t, entry.Amount, -amount)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err)

		//get account
		account2 := result.Account
		require.NotEmpty(t, account2)
		require.Equal(t, account2.ID, account1.ID)

		//check accounts' balance diffrence
		fmt.Println(">> tx: ", account2.Balance)
		diff := account1.Balance - account2.Balance
		require.True(t, diff > 0)

		k := int(diff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	close(errs)
	close(results)
	//check the final updated balances
	updateAccount, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	fmt.Println(">> after: ", updateAccount.Balance)
	require.Equal(t, account1.Balance-float64(m)*amount, updateAccount.Balance)
}

func TestSQLStore_DepositTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance)

	n := 5
	amount := float64(10.00)

	errs := make(chan error)
	results := make(chan BusinessTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.DepositTx(context.Background(), BusinessTxParms{
				AccountID: account1.ID,
				Amount:    amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		//check Entries
		entry := result.Entry
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account1.ID)
		require.Equal(t, entry.Amount, amount)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err)

		//get account
		account2 := result.Account
		require.NotEmpty(t, account2)
		require.Equal(t, account2.ID, account1.ID)

		//check accounts' balance diffrence
		fmt.Println(">> tx: ", account2.Balance)
		diff := account2.Balance - account1.Balance
		require.True(t, diff > 0)

		k := int(diff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final updated balances
	updateAccount, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	fmt.Println(">> after: ", updateAccount.Balance)
	require.Equal(t, account1.Balance+float64(n)*amount, updateAccount.Balance)
}
