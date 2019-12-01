package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func Initialize() {

	poolConfig, err := pgxpool.ParseConfig(os.Getenv("PG_URL"))
	if err != nil {
		log.Fatal("Unable to parse DATABASE_URL", "error", err)
	}

	Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Unable to create connection pool", "error", err)
	}

}
