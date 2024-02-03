package main

import (
	"bank/api"
	db "bank/db/sqlc"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort string = "8080"

func main() {
	database, err := sql.Open("pgx", "host=localhost user=postgres password=password port=5432 database=bank sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err = database.Ping(); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	store := db.NewDBStore(database)
	server := api.NewServer(store)
	server.Serve(fmt.Sprintf(":%s", webPort))
}
