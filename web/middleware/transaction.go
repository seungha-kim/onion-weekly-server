package middleware

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

func Transaction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		pgxPool := c.Get("pgxPool").(*pgxpool.Pool)
		ctx := context.Background()
		tx, err := pgxPool.Begin(ctx)
		if err != nil {
			return
		}
		_, err = tx.Begin(ctx)
		if err != nil {
			return
		}
		c.Set("tx", tx)
		if err = next(c); err != nil {
			err = tx.Rollback(ctx)
			return
		}
		if err = tx.Commit(ctx); err != nil {
			return
		}
		return
	}
}
