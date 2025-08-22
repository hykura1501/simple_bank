package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hykura1501/simple_bank/ulti"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	config, err := ulti.LoadConfig("../../")

	if err != nil {
		log.Fatalf("fail to load the configuration %s", err.Error())
	}
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
