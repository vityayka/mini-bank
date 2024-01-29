package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("pgx", "host=localhost user=postgres password=password port=5432 database=bank sslmode=disable")
	defer testDB.Close()
	if err != nil {
		log.Fatal("couldn't connect to DB", err)
	}
	if err = testDB.Ping(); err != nil {
		log.Fatal("couldn't connect to DB", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
