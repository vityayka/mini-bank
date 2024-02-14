package db

import (
	"bank/utils"
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
	config, err := utils.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
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
