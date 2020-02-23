package middleware

import (
	"context"

	"github.com/onion-studio/onion-weekly/config"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

func Transaction(appConf config.AppConf, pgxPool *pgxpool.Pool) echo.MiddlewareFunc {
	if appConf.Test {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
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
				err = next(c)
				_ = tx.Rollback(ctx)
				return err
			}
		}
	} else {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
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
					_ = tx.Rollback(ctx)
					return
				}
				if err = tx.Commit(ctx); err != nil {
					return
				}
				return
			}
		}
	}

}
