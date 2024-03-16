package main

import (
	db "bank/db/sqlc"
	"bank/gapi"
	"bank/pb"
	"bank/utils"
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

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
	go runGatewayServer(config, store)
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

//go:embed doc/swagger/*
var swaggerFS embed.FS

func runGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	subFS, err := fs.Sub(swaggerFS, "doc/swagger")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(subFS))

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fileServer))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
	}
}
