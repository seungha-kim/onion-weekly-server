package middleware

import (
	"github.com/labstack/echo"
	"github.com/onion-studio/onion-weekly/config"
)

func AppConf(appConf config.AppConf) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("appConf", appConf)
			return next(c)
		}
	}
}
