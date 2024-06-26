package db

import (
	"context"
	"fmt"
)

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *DBStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		swapAccountsToPreventDeadlock(&arg)

		fromAccount, err := queries.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		fmt.Println("fromacc.balance = ", fromAccount.Balance)
		if fromAccount.Balance < arg.Amount {
			result.FromAccount = fromAccount
			return fmt.Errorf("insufficient funds")
		}
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, result.ToAccount, err = addMoney(queries, ctx, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("tx error: %v", err)
	}

	return result, err
}

func swapAccountsToPreventDeadlock(arg *TransferTxParams) {
	if arg.FromAccountID > arg.ToAccountID {
		arg.FromAccountID, arg.ToAccountID = arg.ToAccountID, arg.FromAccountID
		arg.Amount = -arg.Amount
	}
}

func addMoney(
	queries *Queries,
	ctx context.Context,
	fromAccountID int64,
	amount1 int64,
	toAccountID int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = queries.AddBalanceToAccount(ctx, AddBalanceToAccountParams{
		ID:     fromAccountID,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = queries.AddBalanceToAccount(ctx, AddBalanceToAccountParams{
		ID:     toAccountID,
		Amount: amount2,
	})

	return
}
