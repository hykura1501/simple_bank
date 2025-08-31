package main

import (
	"context"
	"log"

	"github.com/hykura1501/simple_bank/api"
	db "github.com/hykura1501/simple_bank/db/sqlc"
	"github.com/hykura1501/simple_bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
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

	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer(config.ServerAddress)

	if err != nil {
		log.Fatal("error when starting server! ", err)
	}

}
