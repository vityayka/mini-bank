package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(context.Context, TransferTxParams) (TransferTxResult, error)
	CreateUserTX(context.Context, CreateUserTxParams) (CreateUserTxResult, error)
}

type DBStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewDBStore(db *pgxpool.Pool) Store {
	return &DBStore{
		Queries:  New(db),
		connPool: db,
	}
}

func (store *DBStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
