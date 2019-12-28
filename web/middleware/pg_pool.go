package middleware

import (
	"github.com/labstack/echo"
	"github.com/onion-studio/onion-weekly/db"
)

func PgxPool(pgUrl string) echo.MiddlewareFunc {
	pool := db.CreatePool(pgUrl)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("pgxPool", pool)
			return next(c)
		}
	}
}
