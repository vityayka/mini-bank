package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	dsn := "host=localhost user=postgres password=password port=5432 database=bank sslmode=disable"
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("couldn't connect to DB", err)
	}
	if err = conn.Ping(); err != nil {
		log.Fatal("couldn't connect to DB", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
