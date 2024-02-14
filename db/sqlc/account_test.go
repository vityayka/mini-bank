package db

import (
	"bank/utils"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandAccount(t *testing.T) (Account, CreateAccountParams) {
	usr, _ := createRandUser(t)
	return createAccountForUser(t, usr.ID, utils.RandomCurrency())
}

func createAccountForUser(t *testing.T, userID int64, currency string) (Account, CreateAccountParams) {
	args := CreateAccountParams{
		Owner:    utils.RandomName(),
		UserID:   userID,
		Balance:  utils.RandomMoney(),
		Currency: currency,
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	return account, args
}

func TestCreateAccount(t *testing.T) {
	account, createArgs := createRandAccount(t)
	require.NotEmpty(t, account)
	require.Equal(t, createArgs.Owner, account.Owner)
	require.Equal(t, createArgs.Balance, account.Balance)
	require.Equal(t, createArgs.Currency, account.Currency)
	require.Equal(t, createArgs.UserID, account.UserID)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	acc1, _ := createRandAccount(t)
	acc2, err := testQueries.GetUserAccount(context.Background(), GetUserAccountParams{acc1.UserID, acc1.ID})
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.UserID, acc2.UserID)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc1, _ := createRandAccount(t)
	args := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: utils.RandomMoney(),
	}
	acc2, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc2.Balance, args.Balance)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.ID, acc2.ID)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1, _ := createRandAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)

	require.NoError(t, err)

	acc2, err := testQueries.GetUserAccount(context.Background(), GetUserAccountParams{acc1.UserID, acc1.ID})

	require.Empty(t, acc2)
	require.ErrorIs(t, sql.ErrNoRows, err)
}

func TestListAccounts(t *testing.T) {
	user, _ := createRandUser(t)

	currencies := utils.GetSupportedCurrencies()
	for _, currency := range currencies {
		createAccountForUser(t, user.ID, currency)
	}

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, accounts, len(currencies))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
