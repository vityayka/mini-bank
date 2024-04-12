package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User `json:"user"`
}

func (store *DBStore) CreateUserTX(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := store.execTx(ctx, func(queries *Queries) error {
		user, err := queries.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			var pgError *pgconn.PgError
			if errors.As(err, &pgError) && pgError.Code == "23505" {
				return ErrUserAlreadyExists
			}
			return err
		}

		result.User = user
		return arg.AfterCreate(user)
	})

	if err != nil {
		return result, err
	}

	return result, nil
}
