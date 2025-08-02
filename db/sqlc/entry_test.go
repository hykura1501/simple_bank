package db

import (
	"context"
	"testing"

	"github.com/hykura1501/simple_bank/ulti"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	acc := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    ulti.RandomMoney(),
	}

	en, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, en)

	require.Equal(t, arg.AccountID, en.AccountID)
	require.Equal(t, arg.Amount, en.Amount)

	require.NotZero(t, en.ID)
	require.NotZero(t, en.CreateAt)

	return en
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	en1 := createRandomEntry(t)
	en2, err := testQueries.GetEntry(context.Background(), en1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, en2)

	require.Equal(t, en1.ID, en2.ID)
	require.Equal(t, en1.AccountID, en2.AccountID)
	require.Equal(t, en1.Amount, en2.Amount)
	require.Equal(t, en1.CreateAt, en2.CreateAt)
}

func TestUpdateEntry(t *testing.T) {
	en1 := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID:     en1.ID,
		Amount: ulti.RandomMoney(),
	}

	en2, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, en2)

	require.Equal(t, en1.ID, en2.ID)
	require.Equal(t, en1.AccountID, en2.AccountID)
	require.Equal(t, arg.Amount, en2.Amount)
	require.Equal(t, en1.CreateAt, en2.CreateAt)
}

func TestDeleteEntry(t *testing.T) {
	en1 := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), en1.ID)
	require.NoError(t, err)

	en2, err := testQueries.GetEntry(context.Background(), en1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, en2)
}

func TestListEntry(t *testing.T) {
	n := 10
	limit := 5
	for range n {
		createRandomEntry(t)
	}
	arg := ListEntriesParams{
		Limit:  int32(limit),
		Offset: 5,
	}

	ens, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, ens, limit)

	for _, en := range ens {
		require.NotEmpty(t, en)
	}
}
