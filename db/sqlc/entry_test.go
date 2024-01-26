package db

import (
	"bank/utils"
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func newCreateEntryParams(account Account) CreateEntryParams {
	return CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}
}

func createNewEntry(t *testing.T, args CreateEntryParams) Entry {
	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	return entry
}

func TestCreateEntry(t *testing.T) {
	account, _ := createRandAccount(t)
	args := newCreateEntryParams(account)
	entry := createNewEntry(t, args)

	require.NotEmpty(t, entry)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)
	require.Equal(t, entry.AccountID, account.ID)
	require.Equal(t, entry.Amount, args.Amount)
}

func TestGetEntry(t *testing.T) {
	account, _ := createRandAccount(t)
	args := newCreateEntryParams(account)
	entry1 := createNewEntry(t, args)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.NotZero(t, entry2.CreatedAt)
	require.NotZero(t, entry2.ID)
	require.Equal(t, entry2.AccountID, entry1.AccountID)
	require.Equal(t, entry2.Amount, entry1.Amount)
}

func TestUpdateEntry(t *testing.T) {
	account, _ := createRandAccount(t)
	args := newCreateEntryParams(account)
	entry := createNewEntry(t, args)

	uArgs := UpdateEntryParams{
		ID:     entry.ID,
		Amount: utils.RandomMoney(),
	}

	uEntry, err := testQueries.UpdateEntry(context.Background(), uArgs)

	require.NoError(t, err)
	require.NotEmpty(t, uEntry)
	require.Equal(t, entry.CreatedAt, uEntry.CreatedAt)
	require.Equal(t, entry.AccountID, uEntry.AccountID)
	require.Equal(t, entry.ID, uEntry.ID)
	require.Equal(t, uEntry.Amount, uArgs.Amount)
}

func TestDeleteEntry(t *testing.T) {
	account, _ := createRandAccount(t)
	args := newCreateEntryParams(account)
	entry := createNewEntry(t, args)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)

	require.NoError(t, err)

	entry1, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.ErrorIs(t, sql.ErrNoRows, err)
	require.Empty(t, entry1)
}

func TestListEntries(t *testing.T) {
	acc, _ := createRandAccount(t)

	for i := 0; i < 10; i++ {
		createNewEntry(t, newCreateEntryParams(acc))
	}

	entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		Limit:  5,
		Offset: 5,
	})

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
