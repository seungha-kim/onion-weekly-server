package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/web"
	m "github.com/onion-studio/onion-weekly/web/middleware"
)

func main() {
	appConf := config.LoadAppConf()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(m.AppConf(appConf))
	e.Use(m.PgxPool(appConf.PgURL))
	e.Static("/", "./public")
	web.RegisterAuthHandlers(e.Group("/auth"))
	web.RegisterHelloHandlers(e.Group("/hello"))
	// FIXME
	e.Debug = appConf.Debug
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", appConf.Port)))
}
