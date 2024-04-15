package main

import (
	"bank/async"
	db "bank/db/sqlc"
	"bank/gapi"
	"bank/pb"
	"bank/utils"
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	database, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().Err(err)
	}
	if err = database.Ping(); err != nil {
		log.Fatal().Err(err)
	}
	defer database.Close()
	migrateDB(config.MigrationURL, config.DBURI)

	store := db.NewDBStore(database)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddr,
	}
	taskDistributor := async.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)

	go runGatewayServer(config, store, taskDistributor)
	startGRPCerver(config, store, taskDistributor)
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := async.NewRedisTaskProcessor(redisOpt, store)
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err)
	}
}

func migrateDB(migrationURL, dbURI string) {
	log.Info().Msg(dbURI)
	log.Info().Msg(migrationURL)
	migrator, err := migrate.New(migrationURL, dbURI)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize the DB migrator")
	}
	if err = migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msgf("0 new migrations have run \n")
			return
		} else {
			log.Fatal().Err(err).Msg("failed to run the DB migration")
		}
	}
	log.Info().Msgf("DB migrations ran successfully \n")
}

func startGRPCerver(config utils.Config, store db.Store, taskDistributor async.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	logger := grpc.UnaryInterceptor(gapi.GRPCLogger)
	grpcServer := grpc.NewServer(logger)

	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer listener.Close()
	log.Info().Msgf("starting grpc server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func runGatewayServer(config utils.Config, store db.Store, taskDistributor async.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
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
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fileServer := http.FileServer(http.Dir("doc/swagger"))

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fileServer))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, gapi.HTTPLogger(mux))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}
