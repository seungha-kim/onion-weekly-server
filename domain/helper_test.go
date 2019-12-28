package domain

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RunInTransaction(pgxPool *pgxpool.Pool, f func(pgx.Tx)) {
	ctx := context.Background()
	tx, err := pgxPool.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	f(tx)
}
