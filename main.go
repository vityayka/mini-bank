package main

import (
	db "bank/db/sqlc"
	"bank/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

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
	// server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	populate(store)
	// server.Serve(fmt.Sprintf(":%s", webPort))
}

func populate(store db.Store) {
	start, err := time.Parse(time.DateOnly, "2005-09-06")
	// now := time.Now()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	for date := start; true; date = date.Add(time.Hour) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(date time.Time) {
				username := utils.RandomString(32)
				user, err := store.CreateUser(ctx, db.CreateUserParams{
					Username:       username,
					HashedPassword: utils.RandomString(32),
					FullName:       fmt.Sprintf("%s %s", username, utils.RandomName()),
					Email:          utils.RandomEmail(),
				})

				if err != nil {
					panic(err)
				}

				err = store.CreateAccount(ctx, db.CreateAccountParams{
					UserID:  user.ID,
					Owner:   user.Username,
					Balance: utils.RandomMoney(),
					// Currency:  utils.RandomCurrency(),
					CreatedAt: date,
				})
				if err != nil {
					panic(err)
				}
				fmt.Println("account created: ", date.String())
				wg.Done()
			}(date)
			wg.Wait()
		}
	}
}
