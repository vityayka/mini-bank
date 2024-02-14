package main

import (
	"bank/api"
	db "bank/db/sqlc"
	"bank/utils"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort string = "8080"

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	database, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal(err)
	}
	if err = database.Ping(); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	store := db.NewDBStore(database)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	server.Serve(fmt.Sprintf(":%s", webPort))
}
