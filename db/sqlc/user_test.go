package db

import (
	"SimpleBank/util"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func (q *Queries) CreateAndReturnUser(ctx context.Context, arg CreateUserParams) (result User, err error) {
	err = q.CreateUser(ctx, arg)
	if err != nil {
		return
	}
	result, err = q.GetUser(ctx, arg.Username)
	return
}

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateAndReturnUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.Equal(t, user1.Username, user1.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user1.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
}
