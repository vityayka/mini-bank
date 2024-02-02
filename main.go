package main

import (
	"bank/api"
	db2 "bank/db/sqlc"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort string = "80"

func main() {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=password port=5432 database=bank sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := db2.NewStore(db)
	server := api.NewServer(store)
	server.Serve(fmt.Sprintf(":%s", webPort))
}
