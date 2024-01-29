package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1, _ := createRandAccount(t)
	acc2, _ := createRandAccount(t)

	amount := int64(10)
	cnt := 10

	result := make(chan TransferTxResult)
	errC := make(chan error)

	for i := 0; i < cnt; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			result <- res
			errC <- err
		}()
	}

	for i := 0; i < cnt; i++ {
		res := <-result
		err := <-errC
		require.NoError(t, err)
		require.NotEmpty(t, res)

		//check entry

		require.Equal(t, amount, res.ToEntry.Amount)
		require.Equal(t, -res.FromEntry.Amount, res.ToEntry.Amount)
		require.Equal(t, res.FromEntry.AccountID, acc1.ID)
		require.Equal(t, res.ToEntry.AccountID, acc2.ID)
		require.NotEmpty(t, res.ToEntry.CreatedAt)
		require.NotEmpty(t, res.FromEntry.CreatedAt)
		require.NotEmpty(t, res.ToEntry.ID)
		require.NotEmpty(t, res.FromEntry.ID)

		_, err = testQueries.GetEntry(context.Background(), res.FromEntry.ID)
		require.NoError(t, err)

		_, err = testQueries.GetEntry(context.Background(), res.ToEntry.ID)
		require.NoError(t, err)

		// check transfer
		_, err = testQueries.GetTransfer(context.Background(), res.Transfer.ID)
		require.NoError(t, err)

		require.Equal(t, res.Transfer.Amount, res.ToEntry.Amount)
		require.Equal(t, res.Transfer.Amount, -res.FromEntry.Amount)
		require.Equal(t, res.Transfer.ToAccountID, acc2.ID)
		require.Equal(t, res.Transfer.FromAccountID, acc1.ID)
		require.NotEmpty(t, res.Transfer.CreatedAt)
		require.NotEmpty(t, res.Transfer.ID)

		//check accounts
		fromAccount := res.FromAccount
		toAccount := res.ToAccount

		require.NotEmpty(t, fromAccount)
		require.NotEmpty(t, toAccount)
		require.Equal(t, fromAccount.ID, acc1.ID)
		require.Equal(t, toAccount.ID, acc2.ID)

		//check balances
		require.Equal(t, fromAccount.Balance, acc1.Balance-amount*int64(i+1))
		require.Equal(t, toAccount.Balance, acc2.Balance+amount*int64(i+1))
	}
}

func TestCreateTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1, _ := createRandAccount(t)
	acc2, _ := createRandAccount(t)

	amount := int64(10)
	cnt := 100

	errC := make(chan error)

	for i := 0; i < cnt; i++ {
		fromAccountId, toAccountId := acc1.ID, acc2.ID
		if i%2 == 0 {
			fromAccountId, toAccountId = acc2.ID, acc1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errC <- err
		}()
	}

	for i := 0; i < cnt; i++ {
		err := <-errC
		require.NoError(t, err)
	}

	fromAccount, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	toAccount, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	//check balances
	require.Equal(t, fromAccount.Balance, acc1.Balance)
	require.Equal(t, toAccount.Balance, acc2.Balance)
}
