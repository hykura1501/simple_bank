package db

import (
	"context"
	"testing"

	"github.com/hykura1501/simple_bank/ulti"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        ulti.RandomMoney(),
	}

	trans, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trans)

	require.Equal(t, arg.FromAccountID, trans.FromAccountID)
	require.Equal(t, arg.ToAccountID, trans.ToAccountID)
	require.Equal(t, arg.Amount, trans.Amount)

	require.NotZero(t, trans.ID)
	require.NotZero(t, trans.CreateAt)

	return trans
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	trans1 := createRandomTransfer(t)
	trans2, err := testQueries.GetTransfer(context.Background(), trans1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, trans2)

	require.Equal(t, trans1.ID, trans2.ID)
	require.Equal(t, trans1.FromAccountID, trans2.FromAccountID)
	require.Equal(t, trans1.ToAccountID, trans2.ToAccountID)
	require.Equal(t, trans1.Amount, trans2.Amount)
	require.Equal(t, trans1.CreateAt, trans2.CreateAt)
}

func TestUpdateTransfer(t *testing.T) {
	trans1 := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     trans1.ID,
		Amount: ulti.RandomMoney(),
	}

	trans2, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, trans2)

	require.Equal(t, trans1.ID, trans2.ID)
	require.Equal(t, trans1.FromAccountID, trans2.FromAccountID)
	require.Equal(t, trans1.ToAccountID, trans2.ToAccountID)
	require.Equal(t, arg.Amount, trans2.Amount)
	require.Equal(t, trans1.CreateAt, trans2.CreateAt)
}

func TestDeleteTransfer(t *testing.T) {
	trans1 := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), trans1.ID)
	require.NoError(t, err)

	trans2, err := testQueries.GetTransfer(context.Background(), trans1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, trans2)
}

func TestListTransfers(t *testing.T) {
	n := 10
	limit := 5
	for range n {
		createRandomTransfer(t)
	}
	arg := ListTransfersParams{
		Limit:  int32(limit),
		Offset: 5,
	}

	trans, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, trans, limit)

	for _, tran := range trans {
		require.NotEmpty(t, tran)
	}
}
