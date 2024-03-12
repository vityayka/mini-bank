package main

import (
	"bank/api"
	db "bank/db/sqlc"
	"bank/gapi"
	"bank/pb"
	"bank/utils"
	"database/sql"
	"log"
	"net"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	startGRPCerver(config, store)
}

func startGRPCerver(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("starting grpc server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

func startHTTPServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	server.Serve(config.HTTPServerAddress)
}
