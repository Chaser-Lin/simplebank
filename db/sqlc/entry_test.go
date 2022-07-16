package db

import (
	"SimpleBank/util"
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

func createRandomEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: util.RandomID(),
		Amount:    util.RandomAmount(),
	}
	entry, err := testQueries.CreateAndReturnEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestQueries_CreateEntries(t *testing.T) {
	createRandomEntry(t)
}

func TestQueries_GetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
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
	arg := GetEntryByAccountParams{
		AccountID: util.RandomID(),
		Limit:     10,
		Offset:    0,
	}
	entries, err := testQueries.GetEntryByAccount(context.Background(), arg)
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
	var id int64 = 3
	err := testQueries.DeleteEntry(context.Background(), id)
	require.NoError(t, err)
	entry, err := testQueries.GetEntry(context.Background(), id)
	require.Error(t, err)
	require.Empty(t, entry)
}
