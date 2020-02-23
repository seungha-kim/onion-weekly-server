package web

import (
	"github.com/labstack/echo/v4"
)

// Hello represents hello world
type Hello struct {
	Hello string `json:"hello"`
	World string `json:"world"`
}

type helloHandler struct{}

func NewHelloHandler() *helloHandler {
	return &helloHandler{}
}

func (h *helloHandler) Register(g *echo.Group) {
	g.GET("", h.handleHello)
}

func (h *helloHandler) handleHello(c echo.Context) error {
	hello := &Hello{
		Hello: "Go Programming",
		World: "Fun!",
	}
	return c.JSON(200, hello)
}
