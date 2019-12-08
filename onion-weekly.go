package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/onion-studio/onion-weekly/db"
	"github.com/onion-studio/onion-weekly/web"
	"log"
)

// Config represent app-wide configuration
type Config struct {
	pgURL string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load .env")
	}

	db.Initialize()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./public")
	web.RegisterAuthHandlers(e.Group("/auth"))
	web.RegisterHelloHandlers(e.Group("/hello"))
	// FIXME
	e.Debug = true
	e.Logger.Fatal(e.Start(":1323"))
}
