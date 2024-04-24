package db

import (
	"bank/utils"
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var testStore Store

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, config.DBURI)
	if err != nil {
		log.Fatal("couldn't connect to DB", err)
	}
	if err = connPool.Ping(ctx); err != nil {
		log.Fatal("couldn't connect to DB", err)
	}

	testStore = NewDBStore(connPool)

	os.Exit(m.Run())
}
