package main

import (
	"context"
	"log"

	"github.com/hykura1501/simple_bank/api"
	db "github.com/hykura1501/simple_bank/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource      = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("Failed to connect db: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.StartServer(serverAddress)

	if err != nil {
		log.Fatal("error when starting server! ", err)
	}

}
