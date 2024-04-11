package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(context.Context, TransferTxParams) (TransferTxResult, error)
	CreateUserTX(context.Context, CreateUserTxParams) (CreateUserTxResult, error)
}

type DBStore struct {
	*Queries
	db *sql.DB
}

func NewDBStore(db *sql.DB) Store {
	return &DBStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *DBStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
