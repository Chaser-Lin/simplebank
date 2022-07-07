package db

import (
	"SimpleBank/db/util"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

//func  (q *Queries) CreateAndReturnEntry(ctx context.Context, arg CreateEntryParams) (result Entry, err error) {
//	err = q.CreateEntry(ctx, arg)
//	if err != nil {
//		return
//	}
//	result, err = q.GetLastEntry(ctx)
//	return
//}

func TestQueries_CreateEntries(t *testing.T) {
	arg := CreateEntryParams{
		AccountID: util.RandomID(),
		Amount:    util.RandomAmount(),
	}
	account, err := testQueries.CreateAndReturnEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, account.AccountID, arg.AccountID)
	require.Equal(t, account.Amount, arg.Amount)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestQueries_GetEntry(t *testing.T) {
	var id int64 = 1
	entry, err := testQueries.GetEntry(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, entry.ID, id)
	require.Equal(t, entry.AccountID, int64(1))
	require.Equal(t, entry.Amount, float64(-100))
	require.NotZero(t, entry.CreatedAt)
}

func TestQueries_ListEntries(t *testing.T) {
	arg := ListEntriesParams{
		Limit:  5,
		Offset: 3,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	for _, entry1 := range entries {
		fmt.Println(entry1)
		entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
		require.NoError(t, err)
		require.Equal(t, entry1.ID, entry2.ID)
		require.Equal(t, entry1.AccountID, entry2.AccountID)
		require.Equal(t, entry1.Amount, entry2.Amount)
		require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	}
}

func TestQueries_GetEntryByAccount(t *testing.T) {
	var accountId int64 = 5
	entries, err := testQueries.GetEntryByAccount(context.Background(), accountId)
	require.NoError(t, err)
	for _, entry1 := range entries {
		fmt.Println(entry1)
		entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
		require.NoError(t, err)
		require.Equal(t, entry1.ID, entry2.ID)
		require.Equal(t, entry1.AccountID, entry2.AccountID)
		require.Equal(t, entry1.Amount, entry2.Amount)
		require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	}
}

func TestQueries_UpdateEntry(t *testing.T) {
	arg := UpdateEntryParams{
		Amount: -100.00,
		ID:     1,
	}
	err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	entry, err := testQueries.GetEntry(context.Background(), arg.ID)
	require.NoError(t, err)
	require.Equal(t, entry.ID, arg.ID)
	require.Equal(t, entry.Amount, arg.Amount)
}

func TestQueries_DeleteEntry(t *testing.T) {
	var id int64 = 11
	err := testQueries.DeleteEntry(context.Background(), id)
	require.NoError(t, err)
	entry, err := testQueries.GetEntry(context.Background(), id)
	require.Error(t, err)
	require.Empty(t, entry)
}
