package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreatePool(pgUrl string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		log.Fatal("Unable to parse DATABASE_URL", "error", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Unable to create connection pool", "error", err)
	}
	return pool
}

func RollbackForTest(pgxPool *pgxpool.Pool, f func(pgx.Tx)) {
	ctx := context.Background()
	tx, err := pgxPool.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	f(tx)
}
