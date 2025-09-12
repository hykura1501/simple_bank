package db

import (
	"context"
	"testing"
	"time"

	"github.com/hykura1501/simple_bank/util"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.Time.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt.Time, user2.PasswordChangedAt.Time, time.Second)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestUpdateUserFullNameOnly(t *testing.T) {
	user := createRandomUser(t)
	newValue := util.RandomOwner()

	arg := UpdateUserParams{
		FullName: &newValue,
		Username: user.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, updatedUser.FullName, user.FullName)
	require.Equal(t, updatedUser.FullName, newValue)

	require.Equal(t, updatedUser.Email, user.Email)
	require.Equal(t, updatedUser.Username, user.Username)
	require.Equal(t, updatedUser.HashedPassword, user.HashedPassword)
	require.WithinDuration(t, updatedUser.PasswordChangedAt.Time, user.PasswordChangedAt.Time, time.Second)
	require.WithinDuration(t, updatedUser.CreatedAt.Time, user.CreatedAt.Time, time.Second)
}

func TestUpdateUserEmailOnly(t *testing.T) {
	user := createRandomUser(t)
	newValue := util.RandomEmail()

	arg := UpdateUserParams{
		Email:    &newValue,
		Username: user.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, updatedUser.Email, user.Email)
	require.Equal(t, updatedUser.Email, newValue)

	require.Equal(t, updatedUser.FullName, user.FullName)
	require.Equal(t, updatedUser.Username, user.Username)
	require.Equal(t, updatedUser.HashedPassword, user.HashedPassword)
	require.WithinDuration(t, updatedUser.PasswordChangedAt.Time, user.PasswordChangedAt.Time, time.Second)
	require.WithinDuration(t, updatedUser.CreatedAt.Time, user.CreatedAt.Time, time.Second)
}

func TestUpdateUserPasswordOnly(t *testing.T) {
	user := createRandomUser(t)
	newValue := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newValue)
	require.NoError(t, err)
	passwordChangedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}

	arg := UpdateUserParams{
		HashedPassword:    &hashedPassword,
		Username:          user.Username,
		PasswordChangedAt: passwordChangedAt,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, updatedUser.HashedPassword, user.HashedPassword)
	require.Equal(t, updatedUser.HashedPassword, hashedPassword)

	require.Equal(t, updatedUser.FullName, user.FullName)
	require.Equal(t, updatedUser.Username, user.Username)
	require.Equal(t, updatedUser.Email, user.Email)

	require.NotEqual(t, updatedUser.PasswordChangedAt, user.PasswordChangedAt)
	require.WithinDuration(t, updatedUser.PasswordChangedAt.Time, passwordChangedAt.Time, time.Second)
	require.WithinDuration(t, updatedUser.CreatedAt.Time, user.CreatedAt.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	password := util.RandomString(6)
	email := util.RandomEmail()
	fullName := util.RandomOwner()
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	passwordChangedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}

	arg := UpdateUserParams{
		HashedPassword:    &hashedPassword,
		Username:          user.Username,
		FullName:          &fullName,
		Email:             &email,
		PasswordChangedAt: passwordChangedAt,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, updatedUser.FullName, user.FullName)
	require.Equal(t, updatedUser.FullName, fullName)

	require.NotEqual(t, updatedUser.HashedPassword, user.HashedPassword)
	require.Equal(t, updatedUser.HashedPassword, hashedPassword)

	require.NotEqual(t, updatedUser.Email, user.Email)
	require.Equal(t, updatedUser.Email, email)

	require.NotEqual(t, updatedUser.PasswordChangedAt, user.PasswordChangedAt)
	require.WithinDuration(t, updatedUser.PasswordChangedAt.Time, passwordChangedAt.Time, time.Second)

	require.Equal(t, updatedUser.Username, user.Username)
	require.WithinDuration(t, updatedUser.CreatedAt.Time, user.CreatedAt.Time, time.Second)

}
