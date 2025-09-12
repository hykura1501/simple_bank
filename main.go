package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hykura1501/simple_bank/api"
	db "github.com/hykura1501/simple_bank/db/sqlc"
	_ "github.com/hykura1501/simple_bank/docs/statik"
	"github.com/hykura1501/simple_bank/gapi"
	"github.com/hykura1501/simple_bank/pb"
	"github.com/hykura1501/simple_bank/util"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal().Msgf("fail to load the configuration %s", err.Error())
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Msgf("Failed to connect db: %s", err)
	}

	// Run db migration
	m, err := migrate.New(
		config.MigrationURL,
		config.DBSource)

	if err != nil {
		log.Fatal().Msgf("Failed to create new migrate instance: %s", err)
	}

	m.Up()
	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("cannot create a grcp server: %s", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot create listener: %s", err)
	}
	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msgf("cannot start gRPC server: %s", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("cannot create a grcp server: %s", err)
	}

	grpcMuxOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(grpcMuxOptions)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msgf("cannot register handler server: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msgf("cannot create statik fs: %s", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot create listener: %s", err)
	}
	log.Info().Msgf("start HTTP server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msgf("cannot start HTTP server: %s", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal().Msgf("cannot create a gin server: %s", err)
	}

	err = server.StartServer(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Msgf("error when starting server! ", err)
	}
}
