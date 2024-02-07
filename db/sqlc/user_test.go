package db

import (
	"bank/utils"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandUser(t *testing.T) (User, CreateUserParams) {
	hashedPassword, err := utils.HashedPassword(utils.RandomString(8))
	require.NoError(t, err)

	args := CreateUserParams{
		Username:       utils.RandomName(),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
		FullName:       fmt.Sprintf("%s %s", utils.RandomName(), utils.RandomName()),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user, args
}

func TestCreateUser(t *testing.T) {
	user, createArgs := createRandUser(t)
	require.NotEmpty(t, user)
	require.Equal(t, createArgs.Username, user.Username)
	require.Equal(t, createArgs.Email, user.Email)
	require.Equal(t, createArgs.HashedPassword, user.HashedPassword)
	require.Equal(t, createArgs.FullName, user.FullName)
	require.Equal(t, createArgs.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	usr1, _ := createRandUser(t)
	usr2, err := testQueries.GetUser(context.Background(), usr1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, usr2)

	require.Equal(t, usr1.Username, usr2.Username)
	require.Equal(t, usr1.FullName, usr2.FullName)
	require.Equal(t, usr1.Email, usr2.Email)
	require.Equal(t, usr1.ID, usr2.ID)
	require.WithinDuration(t, usr1.CreatedAt, usr2.CreatedAt, time.Second)
}
