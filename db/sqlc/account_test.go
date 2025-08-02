package db

import (
	"context"
	"testing"

	"github.com/hykura1501/simple_bank/ulti"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    ulti.RandomOwner(),
		Balance:  ulti.RandomMoney(),
		Currency: ulti.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreateAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.CreateAt, acc2.CreateAt)
}

func TestUpdateAccont(t *testing.T) {
	acc1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: ulti.RandomMoney(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, arg.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.CreateAt, acc2.CreateAt)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListAccount(t *testing.T) {
	n := 10
	limit := 5
	for range n {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  int32(limit),
		Offset: 5,
	}

	accs, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accs, limit)

	for _, acc := range accs {
		require.NotEmpty(t, acc)
	}
}
