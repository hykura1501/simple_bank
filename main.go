package main

import (
	"context"
	"log"
	"net"

	"github.com/hykura1501/simple_bank/api"
	db "github.com/hykura1501/simple_bank/db/sqlc"
	"github.com/hykura1501/simple_bank/gapi"
	"github.com/hykura1501/simple_bank/pb"
	"github.com/hykura1501/simple_bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("fail to load the configuration %s", err.Error())
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect db: ", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("error when starting server! ", err)
	}
}
