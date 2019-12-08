package web

import (
	"github.com/labstack/echo"
)

// Hello represents hello world
type Hello struct {
	Hello string `json:"hello"`
	World string `json:"world"`
}

func RegisterHelloHandlers(g *echo.Group) {
	g.GET("/", handleHello)
}

func handleHello(c echo.Context) error {
	h := &Hello{
		Hello: "Go Programming",
		World: "Fun!",
	}
	return c.JSON(200, h)
}
