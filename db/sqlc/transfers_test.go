package db

import (
	"SimpleBank/db/util"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

//func  (q *Queries) CreateAndReturnTransfer(ctx context.Context, arg CreateTransferParams) (result Transfer, err error) {
//	err = q.CreateTransfer(ctx, arg)
//	if err != nil {
//		return
//	}
//	result, err = q.GetLastTransfer(ctx)
//	return
//}

func createRandomTransfer(t *testing.T) Transfer {
	arg := CreateTransferParams{
		FromAccountID: util.RandomID(),
		ToAccountID:   util.RandomID(),
		Amount:        util.RandomBalance(),
	}
	transfer, err := testQueries.CreateAndReturnTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestQueries_GetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.Equal(t, transfer1.ID, transfer1.ID)
	require.Equal(t, transfer1.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestQueries_GetTransferByFromAccount(t *testing.T) {
	arg := GetTransferByFromAccountParams{
		FromAccountID: util.RandomID(),
		Limit:         10,
		Offset:        0,
	}
	tranfers, err := testQueries.GetTransferByFromAccount(context.Background(), arg)
	require.NoError(t, err)
	for _, tranfer1 := range tranfers {
		//fmt.Println(tranfer1)
		tranfer2, err := testQueries.GetTransfer(context.Background(), tranfer1.ID)
		require.NoError(t, err)
		require.Equal(t, tranfer1.ID, tranfer2.ID)
		require.Equal(t, tranfer1.FromAccountID, tranfer2.FromAccountID)
		require.Equal(t, tranfer1.ToAccountID, tranfer2.ToAccountID)
		require.Equal(t, tranfer1.Amount, tranfer2.Amount)
		require.Equal(t, tranfer1.CreatedAt, tranfer2.CreatedAt)
	}
}
func TestQueries_GetTransferByFromAccountAndToAccount(t *testing.T) {
	arg := GetTransferByFromAccountAndToAccountParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Limit:         10,
		Offset:        0,
	}
	tranfers, err := testQueries.GetTransferByFromAccountAndToAccount(context.Background(), arg)
	require.NoError(t, err)
	for _, tranfer1 := range tranfers {
		//fmt.Println(tranfer1)
		tranfer2, err := testQueries.GetTransfer(context.Background(), tranfer1.ID)
		require.NoError(t, err)
		require.Equal(t, tranfer1.ID, tranfer2.ID)
		require.Equal(t, tranfer1.FromAccountID, tranfer2.FromAccountID)
		require.Equal(t, tranfer1.ToAccountID, tranfer2.ToAccountID)
		require.Equal(t, tranfer1.Amount, tranfer2.Amount)
		require.Equal(t, tranfer1.CreatedAt, tranfer2.CreatedAt)
	}
}

func TestQueries_ListTransfers(t *testing.T) {
	arg := ListTransfersParams{
		Limit:  10,
		Offset: 0,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	for _, transfer1 := range transfers {
		transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
		require.NoError(t, err)
		require.Equal(t, transfer1.ID, transfer2.ID)
		require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
		require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
		require.Equal(t, transfer1.Amount, transfer2.Amount)
		require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
	}
}

func TestQueries_UpdateTransfer(t *testing.T) {
	arg := UpdateTransferParams{
		Amount: 100.00,
		ID:     2,
	}
	err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	transfer, err := testQueries.GetTransfer(context.Background(), arg.ID)
	require.NoError(t, err)
	require.Equal(t, transfer.ID, arg.ID)
	require.Equal(t, transfer.Amount, arg.Amount)
}

func TestQueries_DeleteTransfer(t *testing.T) {
	var id int64 = 3
	err := testQueries.DeleteTransfer(context.Background(), id)
	require.NoError(t, err)
	transfer, err := testQueries.GetTransfer(context.Background(), id)
	require.Error(t, err)
	require.Empty(t, transfer)
}
